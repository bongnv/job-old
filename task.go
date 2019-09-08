package task

import (
	"context"
	"fmt"
)

// Task provides method Wait for a task.
type Task interface {
	// Wait waits the task to finish.
	// It returns error if there is any.
	Wait(ctx context.Context) error
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
	defer func() {
		if r := recover(); r != nil {
			t.err = fmt.Errorf("task: panic recovered with %v", r)
		}
	}()

	if err := waitFor(ctx, t.dependencies); err != nil {
		t.err = err
		return
	}

	if err := t.doer.Do(ctx); err != nil {
		t.err = err
		return
	}
}

func waitFor(ctx context.Context, tasks []Task) error {
	for _, t := range tasks {
		if err := t.Wait(ctx); err != nil {
			return err
		}
	}

	return nil
}
