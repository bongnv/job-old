package task

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_jobImpl(t *testing.T) {
	called := false
	doer := DoFunc(func(_ context.Context) error {
		called = true
		return nil
	})

	j := &jobImpl{}
	_ = j.Run(context.Background(), doer)
	err := j.Wait(context.Background())
	assert.NoError(t, err)
	assert.True(t, called)
}
