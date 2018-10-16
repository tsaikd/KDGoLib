package bucketRateLimiter

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// Limit defines the maximum frequency of some events.
// Limit is represented as number of events per second.
// A zero Limit allows no events.
type Limit = rate.Limit

// Every converts a minimum time interval between events to a Limit.
func Every(interval time.Duration) Limit {
	return rate.Every(interval)
}

func newLimiter(r Limit, b int) *limiterWrap {
	return &limiterWrap{
		Limiter: *rate.NewLimiter(r, b),
	}
}

type limiterWrap struct {
	rate.Limiter

	lastWait time.Time
}

func (t *limiterWrap) WaitN(ctx context.Context, n int) (err error) {
	t.lastWait = time.Now()
	return t.Limiter.WaitN(ctx, n)
}

func (t *limiterWrap) outdated(timeout time.Duration) bool {
	return t.lastWait.Add(timeout).Before(time.Now())
}
