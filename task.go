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

// Run executes a Doer with a list of dependencies.
// It will wait for all dependencies to finish before start Doer.
func Run(ctx context.Context, doer Doer, dependencies ...Task) Task {
	t := &asyncTask{
		done: make(chan struct{}),
	}

	go startTask(ctx, t, doer, dependencies)
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

type asyncTask struct {
	done chan struct{}
	err  error
}

// Wait implements Task interface.
func (t *asyncTask) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.done:
		return t.err
	}
}

func startTask(ctx context.Context, t *asyncTask, doer Doer, dependencies []Task) {
	defer close(t.done)
	defer func() {
		if r := recover(); r != nil {
			t.err = fmt.Errorf("task: panic while executing: %v", r)
		}
	}()

	if err := Wait(ctx, dependencies...); err != nil {
		t.err = err
		return
	}

	if err := doer.Do(ctx); err != nil {
		t.err = err
		return
	}
}
