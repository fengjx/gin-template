package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appLog "gin-template/internal/app/log"
	"gin-template/pkg/errs"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer errs.RecoverFunc(func(rec any, stack *errs.Stack) {
			appLog.ErrorCtx(c.Request.Context(), "panic recovered",
				zap.Any("panic", rec),
				zap.String("stack", fmt.Sprintf("%#v", stack.StackTrace())),
			)
			abortProblem(c, http.StatusInternalServerError, "服务内部错误", "")
		})
		c.Next()
	}
}
