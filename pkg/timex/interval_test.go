package timex

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestSetIntervalRunsUntilContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var count atomic.Int32
	done := make(chan struct{})

	go func() {
		SetInterval(ctx, 10*time.Millisecond, func(_ context.Context) {
			if count.Add(1) >= 3 {
				cancel()
			}
		})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("expected set interval to stop after context cancellation")
	}

	if count.Load() < 3 {
		t.Fatalf("expected callback to run at least 3 times, got %d", count.Load())
	}
}

func TestSetIntervalIgnoresInvalidArguments(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	SetInterval(ctx, 0, func(_ context.Context) {
		t.Fatal("callback should not run when interval is invalid")
	})

	SetInterval(ctx, 10*time.Millisecond, nil)
}
