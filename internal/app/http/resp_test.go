package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"gin-template/internal/app/berr"
	"gin-template/internal/app/config"
	appEnv "gin-template/internal/app/env"
)

func TestAbortWritesUnifiedProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/abort", func(c *gin.Context) {
		c.Set("trace_id", "trace-123")
		Abort(c, berr.ErrUsernameExists.WithDetail("用户名已存在"))
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/abort", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", recorder.Code)
	}

	var problem struct {
		Status  int             `json:"status"`
		Msg     string          `json:"msg"`
		Details string          `json:"details"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &problem); err != nil {
		t.Fatalf("unmarshal problem: %v", err)
	}

	if problem.Msg != "用户名已存在" {
		t.Fatalf("expected msg, got %s", problem.Msg)
	}
	if problem.Details != "用户名已存在" {
		t.Fatalf("expected details, got %s", problem.Details)
	}
	if problem.Status != berr.StatusUsernameExists {
		t.Fatalf("expected business status %d, got %d", berr.StatusUsernameExists, problem.Status)
	}
	if strings.TrimSpace(string(problem.Data)) != "null" {
		t.Fatalf("expected null data, got %s", string(problem.Data))
	}
	if recorder.Header().Get("X-Trace-Id") != "trace-123" {
		t.Fatalf("expected trace id header, got %s", recorder.Header().Get("X-Trace-Id"))
	}
}

func TestRecoveryReturnsUnifiedProblem(t *testing.T) {
	config.ResetForTest()
	appEnv.ResetForTest()
	if err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}

	engine, err := NewEngine()
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	engine.GET("/panic-test", func(c *gin.Context) {
		panic("boom")
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/panic-test", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", recorder.Code)
	}

	var problem struct {
		Status  int             `json:"status"`
		Msg     string          `json:"msg"`
		Details string          `json:"details"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &problem); err != nil {
		t.Fatalf("unmarshal problem: %v", err)
	}

	if problem.Msg != "服务内部错误" {
		t.Fatalf("expected msg 服务内部错误, got %s", problem.Msg)
	}
	if strings.TrimSpace(problem.Details) == "" {
		t.Fatal("expected details in recovery response")
	}
	if strings.TrimSpace(string(problem.Data)) != "null" {
		t.Fatalf("expected null data, got %s", string(problem.Data))
	}
	if strings.TrimSpace(recorder.Header().Get("X-Trace-Id")) == "" {
		t.Fatal("expected trace id in recovery response header")
	}
}
