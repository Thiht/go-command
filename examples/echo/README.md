# echo

**echo** is a basic command showing the usage of a single command with a handler, flags, and positional arguments.

## Build & Run

```sh
go build
```

## Usage

```
$ ./echo
Usage: ./echo [OPTIONS] COMMAND

Example command

Options:
  -verbose
        Enable verbose output

Subcommands:
  echo
```

```
$ ./echo echo -h
Usage: ./echo echo [OPTIONS]

Options:
  -case string
        Case to use (upper, lower)
  -verbose
        Enable verbose output
```

```
$ ./echo echo -verbose -case upper Hello, World!
command echo called with case: upper
HELLO, WORLD!
```
