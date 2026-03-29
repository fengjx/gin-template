package http

import (
	stdhttp "net/http"

	"github.com/gin-gonic/gin"

	"gin-template/internal/app/berr"
)

type responsePayload struct {
	Status  int    `json:"status"`
	Msg     string `json:"msg"`
	Details string `json:"details,omitempty"`
	Data    any    `json:"data"`
}

func OpenAPIErrorHandler(c *gin.Context, err error, status int) {
	Abort(c, berr.FromHTTPStatus(status).WithError(err).WithDetail(err.Error()))
}

func OK(c *gin.Context, body any) {
	c.JSON(stdhttp.StatusOK, responsePayload{
		Status: berr.StatusOK,
		Msg:    "ok",
		Data:   body,
	})
}

// Abort 统一输出统一 JSON 响应。
func Abort(c *gin.Context, err error) {
	berr.Abort(c, err)
}
