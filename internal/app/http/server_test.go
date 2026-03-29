package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gin-template/internal/app/config"
	appEnv "gin-template/internal/app/env"
)

func TestNewEngineServesDocsAndIndex(t *testing.T) {
	config.ResetForTest()
	appEnv.ResetForTest()
	if err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}

	engine, err := NewEngine()
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	docsRecorder := httptest.NewRecorder()
	docsRequest := httptest.NewRequest(http.MethodGet, "/docs", nil)
	engine.ServeHTTP(docsRecorder, docsRequest)
	if docsRecorder.Code != http.StatusOK {
		t.Fatalf("expected /docs 200, got %d", docsRecorder.Code)
	}
	if !strings.Contains(docsRecorder.Body.String(), "SwaggerUIBundle") {
		t.Fatalf("expected swagger ui html, got %s", docsRecorder.Body.String())
	}

	indexRecorder := httptest.NewRecorder()
	indexRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	engine.ServeHTTP(indexRecorder, indexRequest)
	if indexRecorder.Code != http.StatusOK {
		t.Fatalf("expected / 200, got %d", indexRecorder.Code)
	}
}

func TestPprofDisabledReturnsNotFound(t *testing.T) {
	config.ResetForTest()
	appEnv.ResetForTest()
	t.Setenv("APP_SERVER_PPROF_ENABLED", "false")
	if err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}

	engine, err := NewEngine()
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	engine.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected /debug/pprof/ 404 when disabled, got %d", recorder.Code)
	}
}
