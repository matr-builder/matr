package matr

import (
	"context"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"
)

var Version = "v0.1.0"

// ContextKey is used to identify matr values in the context
type ContextKey string

const (
	ctxArgsKey ContextKey = "matr_args"
)

// Matr is the root structure
type Matr struct {
	tasks  map[string]*Task
	onExit func(context.Context, error)
}

// New creates a new Matr struct instance and returns a point to it
func New() *Matr {
	return &Matr{
		tasks: map[string]*Task{},
	}
}

// TaskNames returns a silce of the available task names
func (m *Matr) TaskNames() []string {
	names := []string{}
	for n := range m.tasks {
		names = append(names, n)
	}
	return names
}

// PrintUsage is a helper function to output the usage docs to stdout
func (m *Matr) PrintUsage(cmd string) {
	var err error

	if cmd != "" {
		for _, c := range m.tasks {
			if c.Name == cmd {
				fmt.Println("matr " + cmd + " :\n")
				fmt.Println(c.Doc)
				fmt.Println("")
				return
			}
		}
		err = errors.New("ERROR: no handler found for target \"" + cmd + "\"")
	}

	fmt.Println("\nUsage: matr <opts> [target] args...")

	fmt.Println("\nTargets:")
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	for _, t := range m.tasks {
		if t.Name == "Default" {
			continue
		}
		fmt.Fprintf(tw, "	%s\t%s\n", t.Name, t.Summary)
	}
	tw.Flush()
	fmt.Println(" ")
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}
}

// Handle registers a new task handler with matr. The Handler will then be referenceable by the provided name,
// if a task is named "default" or "" that function will be called if no function name is provided. The
// default function is also a good place to output usage information for the available tasks.
// CallOptions can be used to allow for before and after Handler middleware functions.
func (m *Matr) Handle(task *Task) {
	if task.Name == "" {
		task.Name = "default"
	}
	m.tasks[task.Name] = task
}

// Run will execute the requested task function with the provided context and arguments.
func (m *Matr) Run(ctx context.Context, args ...string) error {
	argsLen := len(args)
	if argsLen > 0 && args[0] == "-h" {
		cmd := ""
		if argsLen > 1 {
			cmd = args[1]
		}
		m.PrintUsage(cmd)
		return nil
	}

	var handlerArgs []string

	taskName := "default"

	if argsLen != 0 {
		taskName = args[0]
	}

	if argsLen > 1 {
		handlerArgs = args[1:]
	}

	ctx = context.WithValue(ctx, ctxArgsKey, handlerArgs)

	task, ok := m.tasks[taskName]
	if !ok {
		m.PrintUsage("")
		return fmt.Errorf("no handler found for target \"%s\"", taskName)
	}

	err := task.Handler(ctx, handlerArgs)
	if m.onExit != nil {
		m.onExit(ctx, err)
	}

	return err
}

// OnExit executes a final function before matr exits
func (m *Matr) OnExit(fn func(ctx context.Context, err error)) {
	m.onExit = fn
	return
}
