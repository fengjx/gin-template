package trace

import (
	"context"

	"github.com/google/uuid"

	"gin-template/internal/app/config"
)

type contextKey string

const (
	traceIDKey contextKey = "trace_id"
	uidKey     contextKey = "uid"
)

func Generate() string {
	return uuid.NewString()
}

func HeaderName() string {
	return config.Get().Trace.HeaderName
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func IDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(traceIDKey).(string); ok {
		return v
	}
	return ""
}

//nolint:revive // 保留兼容旧调用方的函数名。
func TraceIDFromCtx(ctx context.Context) string {
	return IDFromCtx(ctx)
}

func WithUID(ctx context.Context, uid int64) context.Context {
	return context.WithValue(ctx, uidKey, uid)
}

func UIDFromCtx(ctx context.Context) int64 {
	if v, ok := ctx.Value(uidKey).(int64); ok {
		return v
	}
	return 0
}
