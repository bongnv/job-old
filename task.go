package task

import (
	"context"
)

// Task provides method Wait for a task.
type Task interface {
	// Wait waits the task to finish.
	// It returns error if there is any.
	Wait(ctx context.Context) error
}

// Run executes a Doer with a list of dependencies.
// It will wait for all dependencies to finish before start Doer.
func Run(ctx context.Context, doer Doer, dependencies ...Task) Task {
	t := &taskImpl{
		done:         make(chan struct{}),
		dependencies: dependencies,
		doer:         doer,
	}

	go t.start(ctx)
	return t
}

// Wait waits for all tasks to finish. It returns error if there is any from one of those tasks.
func Wait(ctx context.Context, tasks ...Task) error {
	for _, t := range tasks {
		if err := t.Wait(ctx); err != nil {
			return err
		}
	}

	return nil
}

type taskImpl struct {
	dependencies []Task
	doer         Doer

	done chan struct{}
	err  error
}

// Wait implements Task interface.
func (t *taskImpl) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.done:
		return t.err
	}
}

func (t *taskImpl) start(ctx context.Context) {
	defer close(t.done)

	if err := Wait(ctx, t.dependencies...); err != nil {
		t.err = err
		return
	}

	if err := t.doer.Do(ctx); err != nil {
		t.err = err
		return
	}
}
