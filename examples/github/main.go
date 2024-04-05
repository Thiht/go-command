package main

import (
	"context"
	"flag"

	"github.com/Thiht/go-command"
	"github.com/Thiht/go-command/examples/github/handlers"
	"github.com/google/go-github/v56/github"
)

func main() {
	client := github.NewClient(nil)

	root := command.Root().Flags(func(flagSet *flag.FlagSet) {
		flagSet.Bool("verbose", false, "Enable verbose output")
	}).Help("Example command")

	reposCommand := root.SubCommand("repos").Help("Manage GitHub repositories")
	{
		reposCommand.SubCommand("list").Action(handlers.ReposListHandler(client)).Flags(func(flagSet *flag.FlagSet) {
			flagSet.String("user", "", "GitHub user")
		}).Help("List repositories of a GitHub user")
	}

	root.Execute(context.Background())
}
