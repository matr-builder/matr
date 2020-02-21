# matr

[![GoDoc](https://godoc.org/github.com/matr-builder/matr?status.svg)](https://godoc.org/github.com/matr-builder/matr)
![](https://img.shields.io/badge/license-MIT-blue.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/matr-builder/matr)](https://goreportcard.com/report/github.com/matr-builder/matr)

## Install

matr command

```bash
$ go get github.com/matr-builder/matr
```

## Usage

Execute target func

```bash
$ matr [target]
```

List available targets and matr flags

```bash
$ matr -h
Usage: matr <opts> [target] args...

Targets:
   build    Build will builds the project
   run      Run will run the project
   proto    Proto will build the protobuf files into golang files
   test     Test runs all go tests
   bench    Bench runs all the go benchmarks
   docker   Docker builds static go binary then builds a docker image for it
```

Get expanded help for a target

```bash
$ matr -h [target]
```

## Matrfile

A matr file is any regular go file. Matrfiles must be marked with a build target of "matr"
and be a package main file. The default Matrfile is `./Matrfile.go` or `./Matrfile` to avoid
conflicts. A custom Matrfile path can be defined with the `-matrfile` flag (example: `matr -matrfile /somepath/yourfile.go`)

Example Matrfile header

```go
// +build matr

package main
```

## Targets

Any exported function that is `func(context.Context) error` is considered a
matr target. If the function returns an error it will print to stdout and cause the matrfile
to exit with an exit with a non 0 exit code.

Comments on the target function will become documentation accessible by running
`matr` (with no target). This will list all the build targets in this directory with the first
sentence from their docs.

A target may be designated the default target, which is run when the user runs
matr with no target specified. To denote the default, create a function named `Default`.
If no default target is specified, running `matr` with no target will print the list of targets
and docs.

## Dependencies

The helper function `matr.Deps(ctx, ...matr.HandlerFunc)` may be passed a context and any number of
functions (they do not have to be exported), and the Deps function will not return until all
declared dependencies have been run (and any dependencies they have are run).

### Example Dependencies

```go
func Build(ctx context.Context) error {
    err := matr.Deps(ctx, F, G)
    fmt.Println("Build running")
    return err
}

func F(ctx context.Context) error {
    err := matr.Deps(ctx, h)
    fmt.Println("f running")
    return err
}

func G(ctx context.Context) error {
    err := matr.Deps(ctx, F)
    fmt.Println("g running")
    return err
}

func h(ctx context.Context) error {
    fmt.Println("h running")
    return err
}
```
