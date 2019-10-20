package job

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_jobImpl(t *testing.T) {
	t.Run("start", func(t *testing.T) {
		j := New()
		called := false
		task := TaskFunc(func(_ context.Context) error {
			called = true
			return nil
		})

		_ = j.Start(context.Background(), task)
		err := j.Wait(context.Background())
		require.NoError(t, err)
		require.True(t, called)
	})
}
