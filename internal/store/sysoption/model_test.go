package sysoption

import (
	"context"
	"path/filepath"
	"testing"

	"gin-template/internal/app/config"
	"gin-template/internal/app/db"
	appEnv "gin-template/internal/app/env"
)

func TestCreateAndQueryOptionWithTypeAndStatus(t *testing.T) {
	db.ResetForTest()
	config.ResetForTest()
	appEnv.ResetForTest()

	t.Setenv("APP_DATABASE_SQLITE_PATH", filepath.Join(t.TempDir(), "sysoption-store.db"))
	if err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}

	_ = db.Get()

	item := &Model{
		OptionKey:   "site_profile",
		OptionValue: `{"name":"gin-template"}`,
		Description: "站点配置",
		IsPublic:    true,
		Type:        TypeJSON,
		Status:      StatusOnline,
	}
	if err := Create(context.Background(), item); err != nil {
		t.Fatalf("create option: %v", err)
	}

	got, err := ByKey(context.Background(), "site_profile")
	if err != nil {
		t.Fatalf("query option by key: %v", err)
	}
	if got.Type != TypeJSON {
		t.Fatalf("expected type %q, got %q", TypeJSON, got.Type)
	}
	if got.Status != StatusOnline {
		t.Fatalf("expected status %q, got %q", StatusOnline, got.Status)
	}

	publicItems, err := GetPublic(context.Background())
	if err != nil {
		t.Fatalf("query public options: %v", err)
	}
	if len(publicItems) == 0 {
		t.Fatal("expected public options to contain bootstrap or created rows")
	}
}

func TestGetPublicFiltersOfflineOptions(t *testing.T) {
	db.ResetForTest()
	config.ResetForTest()
	appEnv.ResetForTest()

	t.Setenv("APP_DATABASE_SQLITE_PATH", filepath.Join(t.TempDir(), "sysoption-public.db"))
	if err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}

	_ = db.Get()

	if err := Create(context.Background(), &Model{
		OptionKey:   "public_offline",
		OptionValue: "offline",
		Description: "下线公开配置",
		IsPublic:    true,
		Type:        TypeString,
		Status:      StatusOffline,
	}); err != nil {
		t.Fatalf("create offline public option: %v", err)
	}

	items, err := GetPublic(context.Background())
	if err != nil {
		t.Fatalf("query public options: %v", err)
	}
	for _, item := range items {
		if item.OptionKey == "public_offline" {
			t.Fatal("expected offline public option to be filtered out")
		}
	}
}
