package matr

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// Args returns the handler args from the context
func Args(ctx context.Context) []string {
	args, ok := ctx.Value(ctxArgsKey).([]string)
	if !ok {
		return []string{}
	}
	return args
}

// Arg returns the handler arg at the given position from the context
func Arg(ctx context.Context, idx int, defaultStr string) (string, bool) {
	args, ok := ctx.Value(ctxArgsKey).([]string)
	if !ok {
		return defaultStr, false
	}
	if len(args) < idx+1 {
		return defaultStr, false
	}
	return args[idx], true
}

// Deps is a helper function to run the given dependent handlers
func Deps(ctx context.Context, fns ...HandlerFunc) error {
	var err error
	args := Args(ctx)

	for _, fn := range fns {
		err := fn(ctx, args)
		if err != nil {
			return err
		}
	}

	return err
}

// Sh is a helper function for executing shell commands
func Sh(cmdStr string, args ...interface{}) *exec.Cmd {
	cmdStr = fmt.Sprintf(cmdStr, args...)
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd
}
