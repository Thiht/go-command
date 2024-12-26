package command

import (
	"flag"
)

// Lookup returns the value of a flag as a T, or its zero value if the flag doesn't exist.
func Lookup[T any](flagSet *flag.FlagSet, name string) T {
	return flagSet.Lookup(name).Value.(flag.Getter).Get().(T)
}
