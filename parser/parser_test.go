// +build matr

package parser

import "context"

// Example handler for build:js that is used as a subtask
func MultiVarIdentHandler(a, b, c string) error {
	return nil
}

// Example handler for build:js that is used as a subtask
func MultiVarStarIdentHandler(a, b, c *string) error {
	return nil
}

// Example handler for build:js that is used as a subtask
func multiVarSelHandler(a, b, c context.Context) error {
	return nil
}

// Example handler for build:js that is used as a subtask
func multiVarStarSelHandler(a, b, c *context.Context) error {
	return nil
}

// Example handler for build:js that is used as a subtask
func starSelHandler(ctx *context.Context) error {
	return nil
}

// Example handler for build:js that is used as a subtask
func spreadStarSelandler(ctx ...*context.Context) error {
	return nil
}

// Example handler for build:js that is used as a subtask
func spreadSelHandler(ctx ...context.Context) error {
	return nil
}

// Example handler for build:js that is used as a subtask
func spreadIdentHandler(s ...string) error {
	return nil
}

// Example handler for build:js that is used as a subtask
func IdentSpreadIdentHandler(a string, s ...string) error {
	return nil
}

// Example handler for build:js that is used as a subtask
func IdentSpreadIdentMultiReturnHandler(a string, s ...string) (string, int, error) {
	return nil
}
