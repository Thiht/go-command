package command

import (
	"flag"
	"time"
)

// FlagType represents the types of [flag] that implement the [flag.Getter] interface.
type FlagType interface {
	bool | int | int64 | uint | uint64 | string | float64 | time.Duration
}

// Lookup returns the value of a flag as a [FlagType], or its zero value if the flag doesn't exist.
func Lookup[T FlagType](flagSet *flag.FlagSet, name string) T {
	return flagSet.Lookup(name).Value.(flag.Getter).Get().(T)
}
