package timeutil

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestContextSleep(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	start := time.Now()
	ContextSleep(context.Background(), 500*time.Millisecond)
	require.WithinDuration(start.Add(500*time.Millisecond), time.Now(), 50*time.Millisecond)

	start = time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error {
		ContextSleep(ctx, 500*time.Millisecond)
		return nil
	})
	eg.Go(func() error {
		cancel()
		return nil
	})
	require.NoError(eg.Wait())
	require.WithinDuration(start, time.Now(), 50*time.Millisecond)
}
