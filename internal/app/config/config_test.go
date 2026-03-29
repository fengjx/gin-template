package config

import (
	"os"
	"path/filepath"
	"testing"

	appEnv "gin-template/internal/app/env"
	"github.com/spf13/pflag"
)

func TestLoadPriority(t *testing.T) {
	ResetForTest()
	appEnv.ResetForTest()
	t.Cleanup(appEnv.ResetForTest)
	t.Setenv("APP_ENV", "dev")
	t.Setenv("APP_SERVER_HOST", "127.0.0.1")

	root := t.TempDir()
	configDir := filepath.Join(root, "configs")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("mkdir configs: %v", err)
	}
	writeFile(t, filepath.Join(configDir, "config.yaml"), "server:\n  port: 3000\n  cors_allow_origins:\n    - http://localhost:5173\n")
	writeFile(t, filepath.Join(configDir, "config.dev.yaml"), "server:\n  port: 3100\n")
	writeFile(t, filepath.Join(configDir, "config.local.yaml"), "server:\n  port: 3200\n")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	flags := pflag.NewFlagSet("config", pflag.ContinueOnError)
	BindFlags(flags)
	if err := flags.Parse([]string{"--port", "4200"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if err := Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}

	current := Get()
	if current.Server.Port != 4200 {
		t.Fatalf("expected CLI port override 4200, got %d", current.Server.Port)
	}
	if current.Server.Host != "127.0.0.1" {
		t.Fatalf("expected env host override, got %s", current.Server.Host)
	}
	if len(Sources()) != 3 {
		t.Fatalf("expected 3 merged config files, got %v", Sources())
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
