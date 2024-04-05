package command_test

import (
	"flag"
	"strings"
	"testing"
	"time"

	"github.com/Thiht/go-command"
)

type stringList []string

var _ flag.Getter = (*stringList)(nil)

func (sl *stringList) Set(value string) error {
	*sl = strings.Split(value, ",")
	return nil
}

func (sl *stringList) String() string {
	return strings.Join(*sl, ",")
}

func (sl *stringList) Get() any {
	return *sl
}

func TestLookup(t *testing.T) {
	t.Parallel()

	fs := flag.NewFlagSet("test", flag.ExitOnError)
	fs.Bool("bool", false, "")
	fs.Int("int", 0, "")
	fs.Int64("int64", 0, "")
	fs.Uint("uint", 0, "")
	fs.Uint64("uint64", 0, "")
	fs.String("string", "info", "")
	fs.Float64("float64", 0, "")
	fs.Duration("duration", 0, "")

	var strings stringList
	fs.Var(&strings, "strings", "")

	if err := fs.Parse([]string{
		"-bool",
		"-int=1",
		"-int64=1",
		"-uint=1",
		"-uint64=1",
		"-string=debug",
		"-float64=3.14",
		"-duration=1s",
		"-strings=foo,bar,baz",
	}); err != nil {
		t.Fatal(err)
	}

	t.Run("lookup bool", func(t *testing.T) {
		t.Parallel()

		if boolValue := command.Lookup[bool](fs, "bool"); !boolValue {
			t.Errorf("boolValue = %v, want %v", boolValue, true)
		}
	})

	t.Run("lookup int", func(t *testing.T) {
		t.Parallel()

		if intValue := command.Lookup[int](fs, "int"); intValue != 1 {
			t.Errorf("intValue = %v, want %v", intValue, 1)
		}
	})

	t.Run("lookup int64", func(t *testing.T) {
		t.Parallel()

		if int64Value := command.Lookup[int64](fs, "int64"); int64Value != 1 {
			t.Errorf("int64Value = %v, want %v", int64Value, 1)
		}
	})

	t.Run("lookup uint", func(t *testing.T) {
		t.Parallel()

		if uintValue := command.Lookup[uint](fs, "uint"); uintValue != 1 {
			t.Errorf("uintValue = %v, want %v", uintValue, 1)
		}
	})

	t.Run("lookup uint64", func(t *testing.T) {
		t.Parallel()

		if uint64Value := command.Lookup[uint64](fs, "uint64"); uint64Value != 1 {
			t.Errorf("uint64Value = %v, want %v", uint64Value, 1)
		}
	})

	t.Run("lookup string", func(t *testing.T) {
		t.Parallel()

		if stringValue := command.Lookup[string](fs, "string"); stringValue != "debug" {
			t.Errorf("stringValue = %v, want %v", stringValue, "debug")
		}
	})

	t.Run("lookup float64", func(t *testing.T) {
		t.Parallel()

		if float64Value := command.Lookup[float64](fs, "float64"); float64Value != 3.14 {
			t.Errorf("float64Value = %v, want %v", float64Value, 3.14)
		}
	})

	t.Run("lookup duration", func(t *testing.T) {
		t.Parallel()

		if durationValue := command.Lookup[time.Duration](fs, "duration"); durationValue != time.Second {
			t.Errorf("durationValue = %v, want %v", durationValue, time.Second)
		}
	})

	t.Run("lookup custom type stringList", func(t *testing.T) {
		t.Parallel()

		if stringsValue := command.Lookup[stringList](fs, "strings"); len(stringsValue) != 3 || stringsValue[0] != "foo" || stringsValue[1] != "bar" || stringsValue[2] != "baz" {
			t.Errorf("stringsValue = %v, want %v", stringsValue, []string{"foo", "bar", "baz"})
		}
	})

	t.Run("lookup non-existent flag", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if recover() == nil {
				t.Error("lookup of a non-existent flag should panic")
			}
		}()

		if stringValue := command.Lookup[string](fs, "non-existent"); stringValue != "" {
			t.Errorf("stringValue = %v, want %v", stringValue, "")
		}
	})

	t.Run("lookup flag with wrong type", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if recover() == nil {
				t.Error("lookup of a flag with wrong type should panic")
			}
		}()

		if stringValue := command.Lookup[string](fs, "bool"); stringValue != "" {
			t.Errorf("stringValue = %v, want %v", stringValue, "")
		}
	})
}
