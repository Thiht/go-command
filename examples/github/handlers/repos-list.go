package handlers

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/Thiht/go-command"
	"github.com/google/go-github/v56/github"
)

func ReposListHandler(ghClient *github.Client) command.Handler {
	return func(ctx context.Context, flagSet *flag.FlagSet, _ []string) int {
		user := command.Lookup[string](flagSet, "user")
		if user == "" {
			log.Printf("missing required flag: user")
			return 1
		}

		repos, _, err := ghClient.Repositories.List(ctx, user, nil)
		if err != nil {
			log.Printf("failed to list repositories: %v", err)
			return 1
		}

		fmt.Printf("Repositories of %s:\n", user)
		for _, repo := range repos {
			fmt.Printf("- %s\n", *repo.Name)
		}

		return 0
	}
}
