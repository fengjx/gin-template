package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	projectassets "gin-template"
	"gin-template/internal/app/config"
	appEnv "gin-template/internal/app/env"
)

func TestOpenAPIHandlerMatchesEmbeddedSpec(t *testing.T) {
	config.ResetForTest()
	appEnv.ResetForTest()
	if err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}

	expected, err := projectassets.ReadOpenAPI()
	if err != nil {
		t.Fatalf("read embedded spec: %v", err)
	}

	engine, err := NewEngine()
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/openapi/openapi.yaml", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if recorder.Body.String() != string(expected) {
		t.Fatalf("openapi handler body mismatch")
	}
}
