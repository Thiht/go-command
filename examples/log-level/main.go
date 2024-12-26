package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Thiht/go-command"
)

func main() {
	root := command.Root().Flags(func(flagSet *flag.FlagSet) {
		flagSet.String("level", "info", "Minimum level of logs to display")
	}).Middlewares(LevelMiddleware)

	root.SubCommand("info").Action(InfoHandler)
	root.SubCommand("error").Action(ErrorHandler)

	root.Execute(context.Background())
}

func LevelMiddleware(next command.Handler) command.Handler {
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

func InfoHandler(ctx context.Context, _ *flag.FlagSet, args []string) int {
	slog.InfoContext(ctx, strings.Join(args, " "))
	return 0
}

func ErrorHandler(ctx context.Context, _ *flag.FlagSet, args []string) int {
	slog.ErrorContext(ctx, strings.Join(args, " "))
	return 0
}
