package middleware

import (
	"strings"

	"gin-template/internal/app/trace"
	"github.com/gin-gonic/gin"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := parseTraceID(c.GetHeader("traceparent"))
		if traceID == "" {
			traceID = c.GetHeader(trace.HeaderName())
		}
		if traceID == "" {
			traceID = trace.Generate()
		}

		c.Set("trace_id", traceID)
		c.Request = c.Request.WithContext(trace.WithTraceID(c.Request.Context(), traceID))
		c.Writer.Header().Set(trace.HeaderName(), traceID)
		c.Next()
	}
}

func parseTraceID(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.Split(header, "-")
	if len(parts) >= 2 && len(parts[1]) == 32 {
		return parts[1]
	}
	return header
}
