# go-command

[![Go Reference](https://pkg.go.dev/badge/github.com/Thiht/go-command.svg)](https://pkg.go.dev/github.com/Thiht/go-command)

**go-command** is a **lightweight** and **easy to use** library for creating **command lines with commands and subcommands**.

This library is built upon the [`flag`](https://pkg.go.dev/flag) package from the standard library.
The declaration of subcommands is inspired by HTTP routers, and go-command encourages you to define the routing of commands in a single file.

Here's what go-command **will** help you with:

- declaring commands and subcommands
- basic documentation on each command and subcommand (`-h`/`-help`)
- coming later: shell completions

go-command is **not** a framework, so here's where it **won't** help you:

- flag validation beyond what's supported by `flag`
  - note that you can do more than you probably expect if you create custom flag types implementing [`flag.Getter`](https://pkg.go.dev/flag#Getter)
- positional arguments validation, you can do it yourself in your handlers
- no specific support for environment variables, you can manage it with `os.Getenv`
- error handling or logging

## Why this library?

1. I wanted to do subcommands with just the standard library but found it hard to do; this is an attempt at making it easier with minimal abstractions
2. I wanted to declare my subcommands in the same way as [net/http.HandleFunc](https://pkg.go.dev/net/http#HandleFunc)
3. I wanted a simpler alternative to [spf13/cobra](https://github.com/spf13/cobra) and [urfave/cli](https://github.com/urfave/cli)
   - go-command doesn't aim at being as full-featured as these

## How to use?

Almost everything that go-command can do is defined by the `Command` interface in [`command.go`](./command.go).

### Define a root command

You can create a new root command with `command.Root()`. This returns a command on which you can bind an action or flags, or create a new subcommand.

```go
root := command.Root()

// Bind an action
root = root.Action(rootHandler)

// Add global flags
root = root.Flags(func(flagSet *flag.FlagSet) {
  flagSet.Bool("verbose", false, "Enable verbose output")
})

// Set a help text
root = root.Help("Example command")

// Or, defined fluently
root := command.Root().Action(rootHandler).Flags(func(flagSet *flag.FlagSet) {
  flagSet.Bool("verbose", false, "Enable verbose output")
}).Help("Example command")
```

### (Optional) Create a middleware

If you want to share behaviours between a command and its subcommands, you can define add middlewares.
Middlewares are simply wrapper functions around action handlers. They can be used for setup and teardown.
As an example, here's a middleware making use of a root `log-level` flag used to set the display level of the logger for the root command and all of its subcommands:

```go
// Define a root command
root := command.Root()

// Add a global log-level flag
root = root.Flags(func(flagSet *flag.FlagSet) {
  flagSet.String("log-level", "info", "Log level (debug, info, warn, error)")
})

// Add a middleware to the root command
root = root.Middlewares(logLevelMiddleware)

// A middleware accepts a command handler and returns a new handler.
// The handler argument needs to be called to execute the following middlewares or the action.
// This lets you add behaviours before or after executing the command actions, and can be used to inject data to the context.
func logLevelMiddleware(next command.Handler) command.Handler {
  return func(ctx context.Context, flagSet *flag.FlagSet, args []string) int {
    switch level := command.Lookup[string](flagSet, "level"); level {
    case "debug":
      slog.SetLogLoggerLevel(slog.LevelDebug)

    case "info":
      slog.SetLogLoggerLevel(slog.LevelInfo)

    case "warn":
      slog.SetLogLoggerLevel(slog.LevelWarn)

    case "error":
      slog.SetLogLoggerLevel(slog.LevelError)

    default:
      fmt.Println("Unknown level")
      return 1
    }

    return next(ctx, flagSet, args)
  }
}
```

### Add some subcommands

You can then add subcommands with `SubCommand`. The root command and subcommands all share the same `Command` interface.

```go
subCommand := root.SubCommand("my-subcommand")

// Bind an action
subCommand = subCommand.Action(subCommandHandler)

// Add global flags
subCommand = subCommand.Flags(func(flagSet *flag.FlagSet) {
  flagSet.String("input-file", "", "Input file location")
})

// Set a help text
subCommand = subCommand.Help("Example subcommand")

// Or, defined fluently
subCommand := root.SubCommand("my-subcommand").Action(subCommandHandler).Flags(func(flagSet *flag.FlagSet) {
  flagSet.String("input-file", "", "Input file location")
}).Help("Example subcommand")
```

Handlers have to satisfy the `Handler` interface. They receive a context, a flag set, and positional arguments.

go-command provides the `Lookup[T]` helper to easily get a flag value, but you can use [`flag.FlagSet.Lookup`](https://pkg.go.dev/flag#FlagSet.Lookup) directly if you prefer to stick with the standard library.

```go
func rootHandler(ctx context.Context, fs *flag.FlagSet, args []string) int {
  verbose := command.Lookup[bool](flagSet, "verbose")

  if err := doStuff(); err != nil {
    if verbose {
      fmt.Println("something went wrong", err)
    }

    return 1
  }

  return 0
}
```

### Execute your command

When your commands are defined, you can call `Execute()` to parse the execution arguments and call the correct subcommand. This function will call `os.Exit()` with the value returned by the handler.

```go
root := command.Root()

// ...

root.Execute(context.Background())
```

## Examples

Full examples are available in [./examples](./examples):

- [echo](./examples/echo): a basic command showing the usage of a single command with a **handler**, **flags**, and **positional arguments**.
- [github](./examples/github): a command showing the usage subcommands with **dependency injection**.
- [log level middleware](./examples/log-level): a command defining a **middleware** to set the log level for all subcommands.

## License

See [LICENSE](./LICENSE.md)
