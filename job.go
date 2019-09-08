package task

import "context"

// Job groups and composes smaller tasks in a simple way.
type Job interface {
	// Run executes a Doer with a list of dependencies.
	// It will wait for all dependencies to finish before start Doer.
	Run(ctx context.Context, doer Doer, dependencies ...Task) Task

	// Wait waits for all tasks to finish. It returns error if there is any from one of those tasks.
	// It doesn't wait for all tasks to finish.
	Wait(ctx context.Context) error
}

// NewJob creates a new Job.
func NewJob() Job {
	return &jobImpl{}
}

type jobImpl struct {
	tasks []Task
}

// Run implements Job interface.
func (j *jobImpl) Run(ctx context.Context, doer Doer, dependencies ...Task) Task {
	t := &taskImpl{
		done:         make(chan struct{}),
		dependencies: dependencies,
		doer:         doer,
	}

	j.tasks = append(j.tasks, t)
	go t.start(ctx)
	return t
}

// Wait implements Job interface.
func (j *jobImpl) Wait(ctx context.Context) error {
	return waitFor(ctx, j.tasks)
}
