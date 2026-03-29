package cache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCacheSetGetAndDel(t *testing.T) {
	c := New[string, string](Options[string, string]{})

	c.Set("notice", "hello")

	value, ok := c.Get("notice")
	if !ok {
		t.Fatal("expected cache hit after set")
	}
	if value != "hello" {
		t.Fatalf("expected value hello, got %s", value)
	}

	c.Del("notice")

	if _, ok := c.Get("notice"); ok {
		t.Fatal("expected cache miss after delete")
	}
}

func TestCacheGetAllReturnsCopy(t *testing.T) {
	c := New[string, int](Options[string, int]{KeyPrefix: "opt."})
	c.Set("notice", 1)

	items := c.GetAll()
	if got := items["notice"]; got != 1 {
		t.Fatalf("expected exposed key notice=1, got %d", got)
	}

	items["notice"] = 2

	value, ok := c.Get("notice")
	if !ok {
		t.Fatal("expected cache hit after reading cloned map")
	}
	if value != 1 {
		t.Fatalf("expected original cache value to remain 1, got %d", value)
	}
}

func TestCacheRefreshReplacesSnapshotAndFiltersPrefix(t *testing.T) {
	ctx := context.Background()

	var current atomic.Pointer[map[string]int]
	first := map[string]int{
		"opt.notice": 1,
		"sys.about":  2,
	}
	current.Store(&first)

	c := New[string, int](Options[string, int]{
		KeyPrefix: "opt.",
		Loader: func(_ context.Context) (map[string]int, error) {
			source := current.Load()
			cloned := make(map[string]int, len(*source))
			for key, value := range *source {
				cloned[key] = value
			}
			return cloned, nil
		},
	})

	if err := c.Refresh(ctx); err != nil {
		t.Fatalf("refresh cache: %v", err)
	}

	if value, ok := c.Get("notice"); !ok || value != 1 {
		t.Fatalf("expected prefixed key to be exposed as business key, got value=%d ok=%v", value, ok)
	}
	if _, ok := c.Get("sys.about"); ok {
		t.Fatal("expected non-prefixed key to be filtered out")
	}

	second := map[string]int{
		"opt.banner": 3,
	}
	current.Store(&second)

	if err := c.Refresh(ctx); err != nil {
		t.Fatalf("refresh cache second snapshot: %v", err)
	}

	if _, ok := c.Get("notice"); ok {
		t.Fatal("expected previous snapshot to be replaced during refresh")
	}
	if value, ok := c.Get("banner"); !ok || value != 3 {
		t.Fatalf("expected refreshed key banner=3, got value=%d ok=%v", value, ok)
	}
}

func TestCacheRefreshAndRunRequireLoader(t *testing.T) {
	c := New[string, string](Options[string, string]{RefreshInterval: time.Millisecond})

	if err := c.Refresh(context.Background()); !errors.Is(err, ErrLoaderNotConfigured) {
		t.Fatalf("expected ErrLoaderNotConfigured from refresh, got %v", err)
	}
	if err := c.Run(context.Background()); !errors.Is(err, ErrLoaderNotConfigured) {
		t.Fatalf("expected ErrLoaderNotConfigured from run, got %v", err)
	}
}

func TestCacheRunIsIdempotentWhenNoLoopConfigured(t *testing.T) {
	var calls atomic.Int32

	c := New[string, string](Options[string, string]{
		Loader: func(_ context.Context) (map[string]string, error) {
			calls.Add(1)
			return map[string]string{"notice": "ready"}, nil
		},
	})

	if err := c.Run(context.Background()); err != nil {
		t.Fatalf("run cache first time: %v", err)
	}
	if err := c.Run(context.Background()); err != nil {
		t.Fatalf("run cache second time: %v", err)
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("expected loader to run once for idempotent run, got %d", got)
	}
}

func TestCacheRunAutoRefreshAndStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		mu    sync.RWMutex
		items = map[string]string{"notice": "old"}
		calls atomic.Int32
	)

	c := New[string, string](Options[string, string]{
		RefreshInterval: 20 * time.Millisecond,
		Loader: func(_ context.Context) (map[string]string, error) {
			calls.Add(1)

			mu.RLock()
			defer mu.RUnlock()

			cloned := make(map[string]string, len(items))
			for key, value := range items {
				cloned[key] = value
			}
			return cloned, nil
		},
	})

	if err := c.Run(ctx); err != nil {
		t.Fatalf("run cache: %v", err)
	}

	if value, ok := c.Get("notice"); !ok || value != "old" {
		t.Fatalf("expected initial value old, got value=%s ok=%v", value, ok)
	}

	mu.Lock()
	items = map[string]string{"notice": "new"}
	mu.Unlock()

	deadline := time.Now().Add(400 * time.Millisecond)
	for time.Now().Before(deadline) {
		if value, ok := c.Get("notice"); ok && value == "new" {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if value, ok := c.Get("notice"); !ok || value != "new" {
		t.Fatalf("expected auto refresh to load new value, got value=%s ok=%v", value, ok)
	}

	cancel()
	time.Sleep(40 * time.Millisecond)

	before := calls.Load()
	time.Sleep(60 * time.Millisecond)
	after := calls.Load()
	if after != before {
		t.Fatalf("expected refresh loop to stop after cancel, got before=%d after=%d", before, after)
	}
}
