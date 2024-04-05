package command

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Handler represents a command function called by [Command.Execute].
// The command flags can be accessed from the FlagSet parameter using [Lookup] or [flag.Lookup].
type Handler func(context.Context, *flag.FlagSet, []string) int

// Command represents any command or subcommand of the application.
type Command interface {
	// SubCommand adds a new subcommand to an existing command.
	SubCommand(string) Command

	// Action sets the action to execute when calling the command.
	Action(Handler) Command

	// Execute runs the command using [os.Args]. It should normally be called on the root command.
	Execute(context.Context)

	// Help sets the help message of a command.
	Help(string) Command

	// Flags is used to declare the flags of a command.
	Flags(func(*flag.FlagSet)) Command
}

type command struct {
	name        string
	help        string
	handler     Handler
	subCommands map[string]*command
	flagSet     *flag.FlagSet
	parent      *command
}

// Root creates a new root command.
func Root() Command {
	command := command{
		name:        os.Args[0],
		subCommands: map[string]*command{},
		flagSet:     flag.CommandLine,
	}

	flag.CommandLine.Usage = command.usage

	return &command
}

func (c *command) SubCommand(name string) Command {
	c.subCommands[name] = &command{
		name:        name,
		subCommands: map[string]*command{},
		flagSet:     flag.NewFlagSet(name, flag.ExitOnError),
		parent:      c,
	}

	c.subCommands[name].flagSet.Usage = c.subCommands[name].usage

	return c.subCommands[name]
}

func (c *command) Action(handler Handler) Command {
	c.handler = handler
	return c
}

func (c *command) Execute(ctx context.Context) {
	command, args := c, os.Args[1:]
	for {
		if err := command.flagSet.Parse(args); err != nil {
			// This should never occur because the flag sets use flag.ExitOnError
			os.Exit(2) // Use 2 to mimick the behavior of flag.ExitOnError
		}

		args = command.flagSet.Args()
		if len(args) == 0 {
			break
		}

		subCommand, ok := command.subCommands[args[0]]
		if !ok {
			break
		}

		command.flagSet.VisitAll(func(f *flag.Flag) {
			subCommand.flagSet.Var(f.Value, f.Name, f.Usage)
		})

		command = subCommand
		args = args[1:]
	}

	if command.handler == nil {
		if len(args) > 0 {
			command.flagSet.SetOutput(os.Stderr)
			fmt.Fprintf(command.flagSet.Output(), "command provided but not defined: %s\n", args[0])
			command.usage()
			os.Exit(2) // Use 2 to mimick the behavior of flag.ExitOnError
		}

		command.usage()
		os.Exit(0)
	}

	os.Exit(command.handler(ctx, command.flagSet, args))
}

func (c *command) Help(help string) Command {
	c.help = help
	return c
}

func (c *command) Flags(flags func(*flag.FlagSet)) Command {
	flags(c.flagSet)
	return c
}

func (c *command) usage() {
	var builder strings.Builder
	output := c.flagSet.Output()
	c.flagSet.SetOutput(&builder)

	fullCommand := []string{c.name}
	for command := c.parent; command != nil; command = command.parent {
		fullCommand = append([]string{command.name}, fullCommand...)
	}

	optionsHint := " [OPTIONS]"

	subCommandHint := ""
	if len(c.subCommands) > 0 {
		subCommandHint = " [COMMAND]"
		if c.handler == nil {
			subCommandHint = " COMMAND"
		}
	}

	builder.WriteString("Usage: ")
	builder.WriteString(strings.Join(fullCommand, " "))
	builder.WriteString(optionsHint)
	builder.WriteString(subCommandHint)
	builder.WriteString("\n")

	if c.help != "" {
		builder.WriteString("\n")
		builder.WriteString(c.help)
		builder.WriteString("\n")
	}

	// TODO: check if there are options to print
	builder.WriteString("\n")
	builder.WriteString("Options:\n")
	c.flagSet.PrintDefaults()

	if len(c.subCommands) > 0 {
		builder.WriteString("\n")
		builder.WriteString("Subcommands:")

		for name, subCommand := range c.subCommands {
			builder.WriteString("\n  ")
			builder.WriteString(name)
			if subCommand.help != "" {
				builder.WriteString("\n\t")
				builder.WriteString(subCommand.help)
			}
		}
	}

	fmt.Fprintln(output, builder.String())
}
