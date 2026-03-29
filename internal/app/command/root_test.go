package command

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gin-template/internal/app/config"
	appEnv "gin-template/internal/app/env"
)

func TestRootPersistentConfigFlagForSubcommands(t *testing.T) {
	t.Run("flag after subcommands", func(t *testing.T) {
		configPath := writeCustomConfig(t)
		assertConfigFlagApplied(t, []string{"config", "verify", "--config", configPath, "--port", "4200"}, configPath)
	})

	t.Run("flag before subcommands", func(t *testing.T) {
		configPath := writeCustomConfig(t)
		assertConfigFlagApplied(t, []string{"--config", configPath, "--port", "4200", "config", "verify"}, configPath)
	})
}

func assertConfigFlagApplied(t *testing.T, args []string, configPath string) {
	t.Helper()
	config.ResetForTest()
	appEnv.ResetForTest()
	t.Cleanup(config.ResetForTest)
	t.Cleanup(appEnv.ResetForTest)

	root := t.TempDir()
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

	if err := os.MkdirAll(filepath.Join(root, "configs"), 0o755); err != nil {
		t.Fatalf("mkdir configs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "configs", "config.yaml"), []byte(strings.Join([]string{
		"app:",
		"  name: base-app",
		"server:",
		"  port: 3000",
		"  cors_allow_origins:",
		"    - http://localhost:5173",
		"database:",
		"  sqlite_path: runtime/data/app.db",
		"auth:",
		"  same_site: lax",
	}, "\n")), 0o644); err != nil {
		t.Fatalf("write base config: %v", err)
	}

	cmd := newRootCommand()
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute command: %v", err)
	}

	current := config.Get()
	if current.Server.Port != 4200 {
		t.Fatalf("expected CLI port override 4200, got %d", current.Server.Port)
	}
	if sources := config.Sources(); len(sources) != 1 || sources[0] != configPath {
		t.Fatalf("expected custom config as source, got %v", sources)
	}
}

func writeCustomConfig(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "custom.yaml")
	content := strings.Join([]string{
		"app:",
		"  name: custom-app",
		"server:",
		"  port: 4567",
		"  cors_allow_origins:",
		"    - http://example.com",
		"database:",
		"  sqlite_path: /tmp/gin-template-custom.db",
		"auth:",
		"  same_site: lax",
	}, "\n")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write custom config: %v", err)
	}
	return path
}
