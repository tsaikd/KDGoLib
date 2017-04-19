package timeutil

import (
	"context"
	"time"
)

// ContextSleep is a sleep function with context
func ContextSleep(ctx context.Context, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	select {
	case <-ctx.Done():
	}
	cancel()
}
