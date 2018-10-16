package bucketRateLimiter

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// limiter := rate.NewLimiter(rate.Every(10*time.Second), 10)
// if err = limiter.Wait(ctx); err != nil {
// 	return
// }

// OptionCleanTimeout bucket limiter will be removed after last wait plus timeout
type OptionCleanTimeout time.Duration

// OptionCleanTicker check bucket limiter should be removed interval
type OptionCleanTicker time.Duration

// New return BucketRateLimiter
func New(
	ctx context.Context,
	limit Limit,
	burst int,
	options ...interface{},
) (limiter *BucketRateLimiter) {
	ctx, cancel := context.WithCancel(ctx)
	cleanTimeout := time.Minute
	tickerInterval := time.Minute

	for _, option := range options {
		switch opt := option.(type) {
		case OptionCleanTimeout:
			cleanTimeout = time.Duration(opt)
		case OptionCleanTicker:
			tickerInterval = time.Duration(opt)
		}
	}

	egCleaner, _ := errgroup.WithContext(ctx)
	limiter = &BucketRateLimiter{
		ctx:            ctx,
		cancel:         cancel,
		limit:          limit,
		burst:          burst,
		cleanTimeout:   cleanTimeout,
		tickerInterval: tickerInterval,

		buckets:   map[string]*limiterWrap{},
		egCleaner: egCleaner,
		mutex:     sync.Mutex{},
	}
	egCleaner.Go(limiter.cleanerLoop)
	return
}

// BucketRateLimiter main struct
type BucketRateLimiter struct {
	ctx            context.Context
	cancel         context.CancelFunc
	limit          Limit
	burst          int
	cleanTimeout   time.Duration
	tickerInterval time.Duration

	buckets   map[string]*limiterWrap
	egCleaner *errgroup.Group
	mutex     sync.Mutex
}

// Close cancel context to stop cleaner loop
func (t *BucketRateLimiter) Close() (err error) {
	t.cancel()
	return t.egCleaner.Wait()
}

// Wait is shorthand for WaitN(ctx, 1).
func (t *BucketRateLimiter) Wait(ctx context.Context, bucket string) (err error) {
	return t.WaitN(ctx, bucket, 1)
}

// WaitN blocks until lim permits n events to happen.
// It returns an error if n exceeds the Limiter's burst size, the Context is
// canceled, or the expected wait time exceeds the Context's Deadline.
// The burst limit is ignored if the rate limit is Inf.
func (t *BucketRateLimiter) WaitN(ctx context.Context, bucket string, n int) (err error) {
	t.mutex.Lock()
	limiter, ok := t.buckets[bucket]
	if !ok {
		limiter = newLimiter(t.limit, t.burst)
		t.buckets[bucket] = limiter
	}
	t.mutex.Unlock()

	return limiter.WaitN(ctx, n)
}

func (t *BucketRateLimiter) cleanerLoop() (err error) {
	ticker := time.NewTicker(t.tickerInterval)

	for {
		select {
		case <-t.ctx.Done():
			return
		case <-ticker.C:
			t.cleanBucket()
		}
	}
}

func (t *BucketRateLimiter) cleanBucket() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for bucket, limiter := range t.buckets {
		if limiter.outdated(t.cleanTimeout) && limiter.AllowN(time.Now(), limiter.Burst()) {
			delete(t.buckets, bucket)
		}
	}
}
