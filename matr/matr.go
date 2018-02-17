package matr

import (
	"context"
	"errors"
)

// ContextKey is used to identify matr values in the context
type ContextKey string

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

// TaskNames returns an []string of the available task names
func (m *Matr) TaskNames() []string {
	names := []string{}
	for n := range m.tasks {
		names = append(names, n)
	}
	return names
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
	var handlerArgs []string

	taskName := "default"

	if len(args) != 0 {
		taskName = args[0]
	}

	if len(args) > 1 {
		handlerArgs = args[1:]
	}

	ctx = context.WithValue(ctx, ContextKey("args"), handlerArgs)

	ctx, err := m.execTask(ctx, taskName)
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

func (m *Matr) execTask(ctx context.Context, name string) (context.Context, error) {
	var err error

	task, ok := m.tasks[name]
	if !ok {
		t, ok := m.tasks["default"]
		if !ok {
			return ctx, errors.New("No Default handler defined")
		}
		return t.Handler(ctx)
	}

	ctx, err = task.Handler(ctx)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}
