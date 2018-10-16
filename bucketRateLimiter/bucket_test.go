package bucketRateLimiter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestBucketRateLimiter(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	eg, ctx := errgroup.WithContext(context.Background())
	burst := 3
	limiter := New(ctx, Every(100*time.Millisecond), burst)
	start := time.Now()

	eg.Go(func() error {
		bucket := "a"
		require.NoError(limiter.Wait(ctx, bucket))
		require.NoError(limiter.Wait(ctx, bucket))
		require.NoError(limiter.Wait(ctx, bucket))
		require.WithinDuration(start, time.Now(), 50*time.Millisecond)
		require.NoError(limiter.Wait(ctx, bucket))
		require.NoError(limiter.Wait(ctx, bucket))
		require.WithinDuration(start.Add(2*100*time.Millisecond), time.Now(), 50*time.Millisecond)
		return nil
	})

	eg.Go(func() error {
		bucket := "b"
		require.NoError(limiter.Wait(ctx, bucket))
		require.NoError(limiter.Wait(ctx, bucket))
		require.NoError(limiter.Wait(ctx, bucket))
		require.WithinDuration(start, time.Now(), 50*time.Millisecond)
		require.NoError(limiter.Wait(ctx, bucket))
		require.NoError(limiter.Wait(ctx, bucket))
		require.WithinDuration(start.Add(2*100*time.Millisecond), time.Now(), 50*time.Millisecond)
		return nil
	})

	eg.Go(func() error {
		bucket := "c"
		require.NoError(limiter.Wait(ctx, bucket))
		require.NoError(limiter.Wait(ctx, bucket))
		require.NoError(limiter.Wait(ctx, bucket))
		require.WithinDuration(start, time.Now(), 50*time.Millisecond)
		require.NoError(limiter.Wait(ctx, bucket))
		require.NoError(limiter.Wait(ctx, bucket))
		require.WithinDuration(start.Add(2*100*time.Millisecond), time.Now(), 50*time.Millisecond)
		return nil
	})

	require.NoError(eg.Wait())
	require.NoError(limiter.Close())
}

func ExampleBucketRateLimiter() {
	ctx := context.Background()
	burst := 2
	limiter := New(ctx, Every(100*time.Millisecond), burst)
	defer func() {
		traceError(limiter.Close())
	}()

	start := time.Now()
	bucketA := "a"
	bucketB := "b"

	traceError(limiter.Wait(ctx, bucketA))
	fmt.Println("do bucket", bucketA)
	traceError(limiter.Wait(ctx, bucketA))
	fmt.Println("do bucket", bucketA)
	traceError(limiter.Wait(ctx, bucketB))
	fmt.Println("do bucket", bucketB)
	traceError(limiter.Wait(ctx, bucketB))
	fmt.Println("do bucket", bucketB)
	traceError(limiter.Wait(ctx, bucketA))
	fmt.Println("do bucket", bucketA)
	traceError(limiter.Wait(ctx, bucketB))
	fmt.Println("do bucket", bucketB)
	fmt.Println(time.Now().Sub(start).Truncate(100 * time.Millisecond))

	// Output:
	// do bucket a
	// do bucket a
	// do bucket b
	// do bucket b
	// do bucket a
	// do bucket b
	// 100ms
}

func traceError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
