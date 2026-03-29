package env

import (
	"os"
	"testing"
)

func TestCurrentDefaultsToDev(t *testing.T) {
	ResetForTest()
	t.Cleanup(ResetForTest)

	if got := Current(); got != Dev {
		t.Fatalf("expected default env dev, got %s", got)
	}
}

func TestCurrentPrefersCLIOverEnv(t *testing.T) {
	ResetForTest()
	t.Cleanup(ResetForTest)
	t.Setenv("APP_ENV", "test")

	oldArgs := os.Args
	os.Args = []string{"gin-template", "serve", "--env", "prod"}
	t.Cleanup(func() {
		os.Args = oldArgs
	})

	if got := Current(); got != Prod {
		t.Fatalf("expected cli env prod, got %s", got)
	}
}

func TestCurrentNormalizesInvalidValueToDev(t *testing.T) {
	ResetForTest()
	t.Cleanup(ResetForTest)
	t.Setenv("APP_ENV", "staging")

	if got := Current(); got != Dev {
		t.Fatalf("expected invalid env fallback to dev, got %s", got)
	}
}
