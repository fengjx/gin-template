package middleware

import (
	"github.com/gin-gonic/gin"

	"gin-template/internal/app/berr"
)

func abortProblem(c *gin.Context, status int, msg, detail string) {
	switch status {
	case 429:
		berr.Abort(c, berr.ErrTooManyRequests.WithDetail(detail))
	default:
		berr.Abort(c, berr.ErrInternalServerError.WithDetail(detail))
	}
}
