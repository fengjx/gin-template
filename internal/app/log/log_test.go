package log

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	appEnv "gin-template/internal/app/env"
	"gin-template/internal/app/trace"
	"go.uber.org/zap/zapcore"
)

func TestFieldsFromContextIncludesTraceIDAndUID(t *testing.T) {
	ctx := trace.WithUID(trace.WithTraceID(context.Background(), "trace-123"), 42)

	fields := FieldsFromContext(ctx)
	if len(fields) < 2 {
		t.Fatalf("expected at least 2 fields, got %d", len(fields))
	}

	encoder := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(encoder)
	}

	if got := encoder.Fields["trace_id"]; got != "trace-123" {
		t.Fatalf("expected trace_id trace-123, got %#v", got)
	}
	if got := encoder.Fields["uid"]; got != int64(42) {
		t.Fatalf("expected uid 42, got %#v", got)
	}
}

func TestResolveSinkProfile(t *testing.T) {
	logger := buildAccessLogger(loadRuntimeConfig(), filepath.Join(t.TempDir(), "access.log"), zapcore.InfoLevel)
	if logger == nil {
		t.Fatalf("expected access logger to be created")
	}
}

func TestResolveAppSinkProfile(t *testing.T) {
	fileProfile := resolveAppSinkProfile(false, true)
	if fileProfile.file == nil || fileProfile.file.format != encoderFormatJSON {
		t.Fatalf("expected app file sink to use json format")
	}
	if fileProfile.stdout != nil {
		t.Fatalf("expected app file sink to disable stdout")
	}

	stdoutProfile := resolveAppSinkProfile(true, false)
	if stdoutProfile.stdout == nil || stdoutProfile.stdout.format != encoderFormatConsole {
		t.Fatalf("expected app stdout sink to use console format")
	}
	if !stdoutProfile.stdout.colorLevel {
		t.Fatalf("expected app stdout sink to keep colored level in dev")
	}
}

func TestLoadRuntimeConfigInDev(t *testing.T) {
	ResetForTest()
	appEnv.ResetForTest()
	t.Cleanup(ResetForTest)
	t.Cleanup(appEnv.ResetForTest)
	t.Setenv("APP_ENV", "dev")

	cfg := loadRuntimeConfig()
	if cfg.AppFilename != "" {
		t.Fatalf("expected dev applog to use stdout, got %s", cfg.AppFilename)
	}
	if cfg.AccessFilename != defaultAccessPath {
		t.Fatalf("expected access log path %s, got %s", defaultAccessPath, cfg.AccessFilename)
	}
	if cfg.AppLevel != zapcore.DebugLevel {
		t.Fatalf("expected dev app level debug, got %s", cfg.AppLevel)
	}
	if !cfg.ConsoleColor {
		t.Fatalf("expected dev console color enabled")
	}
}

func TestLoadRuntimeConfigInProd(t *testing.T) {
	ResetForTest()
	appEnv.ResetForTest()
	t.Cleanup(ResetForTest)
	t.Cleanup(appEnv.ResetForTest)
	t.Setenv("APP_ENV", "prod")

	cfg := loadRuntimeConfig()
	if cfg.AppFilename != defaultAppPath {
		t.Fatalf("expected prod applog path %s, got %s", defaultAppPath, cfg.AppFilename)
	}
	if cfg.AccessFilename != defaultAccessPath {
		t.Fatalf("expected prod access log path %s, got %s", defaultAccessPath, cfg.AccessFilename)
	}
	if cfg.AppLevel != zapcore.InfoLevel {
		t.Fatalf("expected prod app level info, got %s", cfg.AppLevel)
	}
	if cfg.ConsoleColor {
		t.Fatalf("expected prod console color disabled")
	}
}

func TestDatedFilename(t *testing.T) {
	got := datedFilename(filepath.Join("runtime", "logs", "app.log"), time.Date(2026, 3, 22, 10, 0, 0, 0, time.Local))
	want := filepath.Join("runtime", "logs", "app-2026-03-22.log")
	if got != want {
		t.Fatalf("expected dated filename %s, got %s", want, got)
	}
}
