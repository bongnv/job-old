package task

import "context"

// Doer is a wrapper of method Do.
type Doer interface {
	// Do includes logic to be executed.
	// error can be returned if there is any.
	Do(ctx context.Context) error
}

// Dofunc is an implementation of Doer from a function.
type DoFunc func(ctx context.Context) error

// Do implements Doer interface.
func (f DoFunc) Do(ctx context.Context) error {
	return f(ctx)
}
