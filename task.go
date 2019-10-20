package job

import "context"

// Task is the interface of a task.
// It includes method Exec to execute a task.
type Task interface {
	// Exec to execute a task.
	// It returns error if there is any.
	Exec(ctx context.Context) error
}

// TaskFunc is an implementation of a Task from a function.
type TaskFunc func(ctx context.Context) error

// Exec implements Task interface.
func (f TaskFunc) Exec(ctx context.Context) error {
	return f(ctx)
}
