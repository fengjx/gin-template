package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"gin-template/internal/app/config"
)

type counter struct {
	Count   int
	ResetAt time.Time
}

var (
	rateMu   sync.Mutex
	counters = map[string]*counter{}
)

func GlobalRateLimit() gin.HandlerFunc {
	return limit("global", config.Get().RateLimit.Global)
}

func CriticalRateLimit() gin.HandlerFunc {
	return limit("critical", config.Get().RateLimit.Critical)
}

func UploadRateLimit() gin.HandlerFunc {
	return limit("upload", config.Get().RateLimit.Upload)
}

func limit(name string, window config.RateLimitWindow) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Get().RateLimit.Enabled || window.Requests <= 0 || window.WindowSeconds <= 0 {
			c.Next()
			return
		}

		key := fmt.Sprintf("%s:%s:%s", name, c.ClientIP(), c.FullPath())
		now := time.Now()
		rateMu.Lock()
		item, ok := counters[key]
		if !ok || now.After(item.ResetAt) {
			item = &counter{Count: 0, ResetAt: now.Add(time.Duration(window.WindowSeconds) * time.Second)}
			counters[key] = item
		}
		item.Count++
		count := item.Count
		reset := item.ResetAt
		rateMu.Unlock()

		c.Header("X-RateLimit-Reset", reset.Format(time.RFC3339))
		if count > window.Requests {
			abortProblem(c, http.StatusTooManyRequests, "请求过于频繁", "")
			return
		}
		c.Next()
	}
}
