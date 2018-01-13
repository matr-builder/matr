// +build matr

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/matr-builder/matr/matr"
	"github.com/matr-builder/tlkn"
)

// Default is an example of overriding the default handler
func DefaultX(ctx context.Context) (context.Context, error) {
	fmt.Println("Running Custom Default HandlerFunc...")
	Build(ctx)
	return ctx, nil
}

// Build is used as and example handler
func Build(ctx context.Context) (context.Context, error) {
	matr.Deps(ctx, Proto, Test)
	fmt.Println("Building...")

	_, err := tlkn.Bash(ctx, `
		ls -la
		echo $GOPATH
		ls -l
	`)
	return ctx, err
}

// Run is used as and example handler
func Run(ctx context.Context) (context.Context, error) {
	matr.Deps(ctx, Build)
	fmt.Println("Running...")
	for {
	}
}

// notExported will run the project
func notExported(ctx context.Context) (context.Context, error) {
	fmt.Println("NotExported...")
	time.Sleep(1 * time.Second)
	return ctx, nil
}

// Proto will build the protobuf files into golang files
func Proto(ctx context.Context) (context.Context, error) {
	out, err := tlkn.Bash(ctx, "echo \"build some proto file\"")
	fmt.Println(string(out))
	return ctx, err
}

// Test is used as and example handler
func Test(ctx context.Context) (context.Context, error) {
	out, err := tlkn.Bash(ctx, `echo "Run unit tests..."`)
	fmt.Println(string(out))

	time.Sleep(1 * time.Second)
	return ctx, err
}

// Bench is used as and example handler
func Bench(ctx context.Context) (context.Context, error) {
	args := matr.Args(ctx)
	fmt.Println(args)

	out, err := tlkn.Bash(ctx, `echo "Run benchmark......"`)
	fmt.Println(string(out))

	return ctx, err
}

// Docker is used as and example handler
func Docker(ctx context.Context) (context.Context, error) {
	out, err := tlkn.Bash(ctx, `echo "Build some docker file...."`)
	fmt.Println(string(out))
	return ctx, err
}

// DockerCompose is used as and example handler
// This is a multi line comment that should
// show up in the full docs
func DockerCompose(ctx context.Context) (context.Context, error) {
	out, err := tlkn.Bash(ctx, `echo "Build some docker-compose file...."`)
	fmt.Println(string(out))
	return ctx, err
}
