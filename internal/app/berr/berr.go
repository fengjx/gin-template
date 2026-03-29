package berr

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gin-template/internal/app/trace"
)

// BusinessError 表示可直接映射到 HTTP Problem 响应的业务异常。
type BusinessError struct {
	HTTPStatus int
	Status     int
	Msg        string
	Detail     string
	Cause      error
}

func New(httpStatus, status int, msg string) *BusinessError {
	return &BusinessError{
		HTTPStatus: httpStatus,
		Status:     status,
		Msg:        msg,
	}
}

func (e *BusinessError) Error() string {
	if e == nil {
		return ""
	}
	if e.Detail != "" {
		return e.Detail
	}
	if e.Msg != "" {
		return e.Msg
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return http.StatusText(e.HTTPStatus)
}

func (e *BusinessError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func (e *BusinessError) WithDetail(detail string) *BusinessError {
	if e == nil {
		return nil
	}
	clone := *e
	clone.Detail = detail
	return &clone
}

func (e *BusinessError) WithError(err error) *BusinessError {
	if e == nil {
		return nil
	}
	clone := *e
	clone.Cause = err
	return &clone
}

// Normalize 将任意 error 归一化为 BusinessError。
func Normalize(err error, fallbackStatus int) *BusinessError {
	if err == nil {
		return defaultError(fallbackStatus, nil)
	}

	var businessErr *BusinessError
	if errors.As(err, &businessErr) {
		return fillDefaults(businessErr, fallbackStatus)
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound), errors.Is(err, sql.ErrNoRows):
		return fillDefaults(ErrResourceNotFound.WithError(err), fallbackStatus)
	case errors.Is(err, http.ErrNoCookie):
		return fillDefaults(ErrRequireLogin.WithError(err), fallbackStatus)
	default:
		return defaultError(fallbackStatus, err)
	}
}

func fillDefaults(err *BusinessError, fallbackStatus int) *BusinessError {
	if err == nil {
		return defaultError(fallbackStatus, nil)
	}

	clone := *err
	if clone.HTTPStatus == 0 {
		clone.HTTPStatus = normalizedStatus(fallbackStatus)
	}
	if clone.Status == 0 {
		clone.Status = defaultBusinessStatus(clone.HTTPStatus)
	}
	if clone.Msg == "" {
		clone.Msg = defaultMessage(clone.HTTPStatus)
	}
	if clone.Detail == "" {
		clone.Detail = clone.Msg
	}
	return &clone
}

func defaultError(status int, cause error) *BusinessError {
	status = normalizedStatus(status)
	msg := defaultMessage(status)
	return &BusinessError{
		HTTPStatus: status,
		Status:     defaultBusinessStatus(status),
		Msg:        msg,
		Detail:     msg,
		Cause:      cause,
	}
}

func FromHTTPStatus(status int) *BusinessError {
	return defaultError(status, nil)
}

func normalizedStatus(status int) int {
	if status <= 0 {
		return http.StatusInternalServerError
	}
	return status
}

func defaultMessage(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "请求参数无效"
	case http.StatusUnauthorized:
		return "请先登录"
	case http.StatusForbidden:
		return "没有权限"
	case http.StatusNotFound:
		return "资源不存在"
	case http.StatusConflict:
		return "资源冲突"
	case http.StatusTooManyRequests:
		return "请求过于频繁"
	default:
		return "服务内部错误"
	}
}

func defaultBusinessStatus(httpStatus int) int {
	switch normalizedStatus(httpStatus) {
	case http.StatusBadRequest:
		return StatusBadRequest
	case http.StatusUnauthorized:
		return StatusUnauthorized
	case http.StatusForbidden:
		return StatusForbidden
	case http.StatusNotFound:
		return StatusNotFound
	case http.StatusConflict:
		return StatusConflict
	case http.StatusTooManyRequests:
		return StatusTooManyRequests
	default:
		return StatusInternalServerError
	}
}

// Abort 统一输出 JSON 响应。
func Abort(c *gin.Context, err error) {
	businessErr := Normalize(err, http.StatusInternalServerError)
	traceID := trace.TraceIDFromCtx(c)
	if traceID == "" {
		traceID = c.GetString("trace_id")
	}
	if traceID != "" {
		c.Header("X-Trace-Id", traceID)
	}

	c.AbortWithStatusJSON(businessErr.HTTPStatus, gin.H{
		"status":  businessErr.Status,
		"msg":     businessErr.Msg,
		"details": businessErr.Error(),
		"data":    nil,
	})
}
