package job

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_taskExecution(t *testing.T) {
	t.Run("context-timeout", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		mockTask := &taskExecution{
			done: make(chan struct{}),
		}

		cancel()
		err := mockTask.Wait(ctx)
		require.EqualValues(t, context.Canceled, err)
	})

	t.Run("task-complete", func(t *testing.T) {
		called := false
		mockTask := &taskExecution{
			done: make(chan struct{}),
		}

		task := TaskFunc(func(_ context.Context) error {
			called = true
			return nil
		})

		startTask(context.Background(), mockTask, task, nil)
		err := mockTask.Wait(context.Background())
		require.NoError(t, err)
		require.True(t, called)
	})

	t.Run("task-panic", func(t *testing.T) {
		mockTask := &taskExecution{
			done: make(chan struct{}),
		}

		task := TaskFunc(func(_ context.Context) error {
			panic("runtime panic")
		})

		require.NotPanics(t, func() {
			startTask(context.Background(), mockTask, task, nil)
		})
	})

	t.Run("task-err", func(t *testing.T) {
		called := false
		mockTask := &taskExecution{
			done: make(chan struct{}),
		}

		task := TaskFunc(func(_ context.Context) error {
			called = true
			return errors.New("runtime error")
		})

		startTask(context.Background(), mockTask, task, nil)
		err := mockTask.Wait(context.Background())
		require.Error(t, err)
		require.True(t, called)
	})

	t.Run("depepencies-error", func(t *testing.T) {
		called := false
		mockTask := &taskExecution{
			done: make(chan struct{}),
		}

		task := TaskFunc(func(_ context.Context) error {
			called = true
			return nil
		})
		dependencies := []Execution{
			mockErrDependency(),
		}

		startTask(context.Background(), mockTask, task, dependencies)
		err := mockTask.Wait(context.Background())
		require.Error(t, err)
		require.False(t, called)
	})
}

func mockErrDependency() Execution {
	mock := &taskExecution{
		done: make(chan struct{}),
		err:  errors.New("random error"),
	}

	close(mock.done)
	return mock
}
