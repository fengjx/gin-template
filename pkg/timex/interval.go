package timex

import (
	"context"
	"time"

	"gin-template/pkg/errs"
)

// SetInterval 会按固定间隔重复执行任务，行为上类似 JavaScript 的 setInterval。
// 调用会阻塞当前 goroutine，直到 context 结束或 interval 非法时直接返回。
func SetInterval(ctx context.Context, interval time.Duration, fn func(ctx context.Context)) {
	if interval <= 0 || fn == nil {
		return
	}

	f := func(ctx context.Context) {
		defer errs.Recover()
		fn(ctx)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f(ctx)
		}
	}
}
