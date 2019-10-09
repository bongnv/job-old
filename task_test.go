package task

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_asyncTask(t *testing.T) {
	t.Run("context-timeout", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		mockTask := &asyncTask{
			done: make(chan struct{}),
		}

		cancel()
		err := mockTask.Wait(ctx)
		require.EqualValues(t, context.Canceled, err)
	})

	t.Run("task-complete", func(t *testing.T) {
		called := false
		mockTask := &asyncTask{
			done: make(chan struct{}),
		}

		doer := DoFunc(func(_ context.Context) error {
			called = true
			return nil
		})

		startTask(context.Background(), mockTask, doer, nil)
		err := mockTask.Wait(context.Background())
		require.NoError(t, err)
		require.True(t, called)
	})

	t.Run("task-panic", func(t *testing.T) {
		mockTask := &asyncTask{
			done: make(chan struct{}),
		}

		doer := DoFunc(func(_ context.Context) error {
			panic("runtime panic")
		})

		require.NotPanics(t, func() {
			startTask(context.Background(), mockTask, doer, nil)
		})
	})

	t.Run("task-err", func(t *testing.T) {
		called := false
		mockTask := &asyncTask{
			done: make(chan struct{}),
		}

		doer := DoFunc(func(_ context.Context) error {
			called = true
			return errors.New("runtime error")
		})

		startTask(context.Background(), mockTask, doer, nil)
		err := mockTask.Wait(context.Background())
		require.Error(t, err)
		require.True(t, called)
	})

	t.Run("depepencies-error", func(t *testing.T) {
		called := false
		mockTask := &asyncTask{
			done: make(chan struct{}),
		}

		doer := DoFunc(func(_ context.Context) error {
			called = true
			return nil
		})
		dependencies := []Task{
			mockErrDependency(),
		}

		startTask(context.Background(), mockTask, doer, dependencies)
		err := mockTask.Wait(context.Background())
		require.Error(t, err)
		require.False(t, called)
	})
}

func mockErrDependency() Task {
	mock := &asyncTask{
		done: make(chan struct{}),
		err:  errors.New("random error"),
	}

	close(mock.done)
	return mock
}

func Test_Run(t *testing.T) {
	called := false
	doer := DoFunc(func(_ context.Context) error {
		called = true
		return nil
	})

	task1 := Run(context.Background(), doer)
	err := Wait(context.Background(), task1)
	require.NoError(t, err)
	require.True(t, called)
}
