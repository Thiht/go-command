# go-command

**go-command** is a **lightweight** and **easy to use** library for creating **command lines with commands and subcommands**.

This library is built upon the [`flag`](https://pkg.go.dev/flag) package from the standard library.
The declaration of subcommands is inspired by HTTP routers, and go-command encourages you to define the routing of commands in a single file.

Here's what go-command **will** help you do:

- easily declare commands and subcommands
- basic documentation on each command and subcommand (`-h`/`-help`)
- coming later: generating shell completions

go-command is **not** a framework, so here's where it **won't** help you:

- flag validation beyond what's supported by `flag`
  - note that you can do more than you probably expect if you create custom flag types implementing [`flag.Getter`](https://pkg.go.dev/flag#Getter)
- positional arguments validation
- no specific support for environment variables
- error handling
- logging

## Why this library?

1. I wanted to do subcommands with just the standard library but found it hard to do ; this is an attempt at making it easier with minimal abstractions
2. I wanted to declare my subcommands in the same way as [net/http.HandleFunc](https://pkg.go.dev/net/http#HandleFunc)
3. I wanted a simpler alternative to [spf13/cobra](https://github.com/spf13/cobra) and [urfave/cli](https://github.com/urfave/cli)
   - go-command doesn't aim at being as full-featured as these

## Examples

Full examples are available in [./examples](./examples):

- [echo](./examples/echo): a basic command showing the usage of a single command with a handler, flags, and positional arguments.
- [github](./examples/github): a command showing the usage subcommands with dependency injection.

---

```go
package main

import (
	"context"
	"flag"
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

	foobarCommand := root.SubCommand("foobar", nil)
	{
		foobarCommand.SubCommand("create", nil)
		foobarCommand.SubCommand("read", nil)
		foobarCommand.SubCommand("update", nil)
		foobarCommand.SubCommand("delete", nil)
	}

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
```

## License

See [LICENSE](./LICENSE.md)
