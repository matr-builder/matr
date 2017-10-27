package matr

import (
	"context"
)

// Task struct that holds registered handler
type Task struct {
	Name    string
	Handler HandlerFunc
	Summary string
	Doc     string
}

// The HandlerFunc type is an adapter to allow the use of ordinary functions as a matr task Handler.
type HandlerFunc func(c context.Context) (context.Context, error)
