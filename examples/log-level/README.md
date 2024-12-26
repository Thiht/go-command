# log-level

**log-level** is a basic command showing how to make use of go-command's middlewares.

## Build & Run

```sh
go build
```

## Usage

```
$ ./log-level info Hello, world!
2024/12/26 22:57:13 INFO Hello, world!
```

```
$ ./log-level -level=warn info Hello, world!
<no output>
```

```
$ ./log-level -level=warn error Hello, world!
2024/12/26 22:58:07 ERROR Hello, world!
```
