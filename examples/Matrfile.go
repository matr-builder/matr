//go:build matr
// +build matr

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/matr-builder/matr/matr"
)

// Default is an example of overriding the default handler
func Default(ctx context.Context, args []string) error {
	fmt.Println("Running Custom Default HandlerFunc...")
	Build(ctx, args)
	return nil
}

// Build is used as and example handler
func Build(ctx context.Context, args []string) error {
	matr.Deps(ctx, Proto, Test)
	fmt.Println("Building...")

	out, err := matr.Sh(`
		ls -la
		echo $GOPATH
		ls -l
	`).Output()
	os.Stdout.Write(out)
	return err
}

// PrintArgs prints the provided args
func PrintArgs(ctx context.Context, args []string) error {
	fmt.Println("args:", "["+strings.Join(args, ",")+"]")
	return nil
}

// Run is used as and example handler
func Run(ctx context.Context, args []string) error {
	matr.Deps(ctx, Build)
	fmt.Println("Running...")
	for {
	}
}

// notExported will run the project
func notExported(ctx context.Context, args []string) error {
	fmt.Println("NotExported...")
	time.Sleep(1 * time.Second)
	return nil
}

// Proto will build the protobuf files into golang files
func Proto(ctx context.Context, args []string) error {
	err := matr.Sh("echo \"build some proto file\"").Run()
	return err
}

// Test is used as and example handler
func Test(ctx context.Context, args []string) error {
	err := matr.Sh(`echo "Run unit tests..."`).Run()
	time.Sleep(1 * time.Second)
	return err
}

// Bench is used as and example handler
func Bench(ctx context.Context, args []string) error {
	err := matr.Sh(`echo "Run benchmark......"`).Run()
	return err
}

// Docker is used as and example handler
func Docker(ctx context.Context, args []string) error {
	err := matr.Sh(`echo "Build some docker file...."`).Run()
	return err
}

// DockerCompose is used as and example handler
// This is a multi line comment that should
// show up in the full docs
func DockerCompose(ctx context.Context, args []string) error {
	err := matr.Sh(`echo "Build some docker-compose file...."`).Run()
	return err
}
