package matr

import "context"

// Args returns the handler args from the context
func Args(ctx context.Context) []string {
	args, ok := ctx.Value(ContextKey("args")).([]string)
	if !ok {
		return []string{}
	}
	return args
}

// Deps is a helper function to run the given dependent handlers
func Deps(ctx context.Context, fns ...HandlerFunc) (context.Context, error) {
	var err error

	for _, fn := range fns {
		ctx, err = fn(ctx)
		if err != nil {
			return ctx, err
		}
	}

	return ctx, err
}
