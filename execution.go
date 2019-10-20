package job

import (
	"context"
	"fmt"
)

// Execution is a wrapper of method Wait.
type Execution interface {
	// Wait waits for a task to finish.
	// It returns error if there is any.
	Wait(ctx context.Context) error
}

type taskExecution struct {
	done chan struct{}
	err  error
}

// Wait implements Task interface.
func (t *taskExecution) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.done:
		return t.err
	}
}

func startTask(ctx context.Context, t *taskExecution, task Task, dependencies []Execution) {
	defer close(t.done)
	defer func() {
		if r := recover(); r != nil {
			t.err = fmt.Errorf("task: panic while executing: %v", r)
		}
	}()

	if err := waitForExecutions(ctx, dependencies); err != nil {
		t.err = err
		return
	}

	if err := task.Exec(ctx); err != nil {
		t.err = err
		return
	}
}

// wait waits for all executions to finish. It returns error if there is any from one of them.
func waitForExecutions(ctx context.Context, executions []Execution) error {
	for _, t := range executions {
		if err := t.Wait(ctx); err != nil {
			return err
		}
	}

	return nil
}
