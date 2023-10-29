package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/Thiht/go-command"
)

func main() {
	root := command.Root(nil).Flags(func(flagSet *flag.FlagSet) {
		flagSet.Bool("verbose", false, "Enable verbose output")
	}).Help("Example command")

	root.SubCommand("echo", EchoHandler).Flags(func(flagSet *flag.FlagSet) {
		flagSet.String("case", "", "Case to use (upper, lower)")
	})

	root.Execute(context.Background())
}

func EchoHandler(ctx context.Context, flagSet *flag.FlagSet, args []string) int {
	verbose := command.Lookup[bool](flagSet, "verbose")
	textCase := command.Lookup[string](flagSet, "case")

	if verbose {
		fmt.Println("command echo called with case: " + textCase)
	}

	switch textCase {
	case "upper":
		fmt.Println(strings.ToUpper(strings.Join(args, " ")))

	case "lower":
		fmt.Println(strings.ToLower(strings.Join(args, " ")))

	default:
		fmt.Println(strings.Join(args, " "))
	}

	return 0
}
