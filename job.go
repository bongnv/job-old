package job

import (
	"context"
)

// Job is the interface of a job. A job is a group of tasks.
type Job interface {
	// Start to start a job.
	Start(ctx context.Context, task Task, dependencies ...Execution) Execution
	// Wait waits for all tasks started by the job.
	Wait(ctx context.Context) error
}

// New creates a new jobk.
func New() Job {
	return &jobImpl{}
}

type jobImpl struct {
	executions []Execution
}

// Run executes a Doer with a list of dependencies.
// It will wait for all dependencies to finish before start Doer.
func (j *jobImpl) Start(ctx context.Context, task Task, dependencies ...Execution) Execution {
	t := &taskExecution{
		done: make(chan struct{}),
	}

	go startTask(ctx, t, task, dependencies)
	j.executions = append(j.executions, t)
	return t
}

func (j *jobImpl) Wait(ctx context.Context) error {
	return waitForExecutions(ctx, j.executions)
}
