package berr

import (
	"errors"
	"net/http"
	"testing"

	"gin-template/pkg/errs"
	"gorm.io/gorm"
)

func TestBusinessErrorUnwrapAndDetail(t *testing.T) {
	root := errors.New("root")
	err := ErrUsernameExists.WithError(root).WithDetail("用户名已存在")

	if err.Error() != "用户名已存在" {
		t.Fatalf("expected detail as error string, got %q", err.Error())
	}
	if !errors.Is(err, root) {
		t.Fatalf("expected errors.Is to match wrapped cause")
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		fallbackStatus int
		wantHTTPStatus int
		wantStatus     int
		wantMsg        string
		wantDetail     string
	}{
		{
			name:           "business error keeps fields",
			err:            ErrUsernameExists.WithError(errors.New("db")).WithDetail("用户名已存在"),
			fallbackStatus: http.StatusInternalServerError,
			wantHTTPStatus: http.StatusConflict,
			wantStatus:     StatusUsernameExists,
			wantMsg:        "用户名已存在",
			wantDetail:     "用户名已存在",
		},
		{
			name:           "record not found maps to not found",
			err:            errs.Wrap(gorm.ErrRecordNotFound, "query user"),
			fallbackStatus: http.StatusInternalServerError,
			wantHTTPStatus: http.StatusNotFound,
			wantStatus:     StatusNotFound,
			wantMsg:        "资源不存在",
			wantDetail:     "资源不存在",
		},
		{
			name:           "fallback bad request",
			err:            errors.New("invalid format"),
			fallbackStatus: http.StatusBadRequest,
			wantHTTPStatus: http.StatusBadRequest,
			wantStatus:     StatusBadRequest,
			wantMsg:        "请求参数无效",
			wantDetail:     "请求参数无效",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.err, tt.fallbackStatus)
			if got.HTTPStatus != tt.wantHTTPStatus {
				t.Fatalf("expected http status %d, got %d", tt.wantHTTPStatus, got.HTTPStatus)
			}
			if got.Status != tt.wantStatus {
				t.Fatalf("expected business status %d, got %d", tt.wantStatus, got.Status)
			}
			if got.Msg != tt.wantMsg {
				t.Fatalf("expected msg %s, got %s", tt.wantMsg, got.Msg)
			}
			if got.Detail != tt.wantDetail {
				t.Fatalf("expected detail %s, got %s", tt.wantDetail, got.Detail)
			}
		})
	}
}
