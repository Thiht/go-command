# github

**github** is a command showing the usage subcommands with dependency injection.

## Build & Run

```sh
go build
```

## Usage

```
$ ./github
Usage: ./github [OPTIONS] COMMAND

Example command

Options:
  -verbose
        Enable verbose output

Subcommands:
  repos
        Manage GitHub repositories
```

```
$ ./github repos list -user octocat
Repositories of octocat:
- boysenberry-repo-1
- git-consortium
- hello-worId
- Hello-World
- linguist
- octocat.github.io
- Spoon-Knife
- test-repo1
```
