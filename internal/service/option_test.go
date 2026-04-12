package service

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/spf13/pflag"

	"gin-template/internal/app/config"
	"gin-template/internal/app/db"
	appEnv "gin-template/internal/app/env"
	sysoptionStore "gin-template/internal/store/sysoption"
)

type optionJSONValue struct {
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
}

func TestGetOptionStringAndJSON(t *testing.T) {
	ResetOptionForTest()
	db.ResetForTest()
	config.ResetForTest()
	appEnv.ResetForTest()

	t.Setenv("APP_DATABASE_SQLITE_PATH", filepath.Join(t.TempDir(), "option-service.db"))
	if err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}

	_ = db.Get()

	if err := sysoptionStore.Create(context.Background(), &sysoptionStore.Model{
		OptionKey:   "site_profile",
		OptionValue: `{"name":"gin-template","enable":true}`,
		Description: "站点配置",
		IsPublic:    false,
		Type:        sysoptionStore.TypeJSON,
		Status:      sysoptionStore.StatusOnline,
	}); err != nil {
		t.Fatalf("create option: %v", err)
	}

	value, err := GetOptionString(context.Background(), "about")
	if err != nil {
		t.Fatalf("get option string: %v", err)
	}
	if value == "" {
		t.Fatal("expected about option value, got empty string")
	}

	profile, err := GetOptionJSON[optionJSONValue](context.Background(), "site_profile")
	if err != nil {
		t.Fatalf("get option json: %v", err)
	}
	if profile.Name != "gin-template" || !profile.Enable {
		t.Fatalf("unexpected json value: %+v", profile)
	}

	_, err = GetOptionString(context.Background(), "missing_key")
	if !errors.Is(err, ErrOptionNotFound) {
		t.Fatalf("expected ErrOptionNotFound, got %v", err)
	}
}

func TestUpdateOptionRefreshesCacheImmediately(t *testing.T) {
	ResetOptionForTest()
	db.ResetForTest()
	config.ResetForTest()
	appEnv.ResetForTest()

	t.Setenv("APP_DATABASE_SQLITE_PATH", filepath.Join(t.TempDir(), "option-update.db"))
	if err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}

	_ = db.Get()

	if err := RefreshOptions(context.Background()); err != nil {
		t.Fatalf("refresh options: %v", err)
	}

	item, err := UpdateOption(context.Background(), "about", UpdateOptionRequest{
		Value:       "新的关于信息",
		Description: "关于信息",
		IsPublic:    true,
		Type:        sysoptionStore.TypeString,
		Status:      sysoptionStore.StatusOnline,
	})
	if err != nil {
		t.Fatalf("update option: %v", err)
	}
	if item.OptionValue != "新的关于信息" {
		t.Fatalf("expected updated item value, got %s", item.OptionValue)
	}

	value, err := GetOptionString(context.Background(), "about")
	if err != nil {
		t.Fatalf("get option string after update: %v", err)
	}
	if value != "新的关于信息" {
		t.Fatalf("expected refreshed cache value, got %s", value)
	}
}

func TestOptionAutoRefresh(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		mu    sync.RWMutex
		items = []sysoptionStore.Model{
			{OptionKey: "notice", OptionValue: "old", Status: sysoptionStore.StatusOnline},
		}
	)

	svc := newOptionWithLoader(20*time.Millisecond, func(_ context.Context) ([]sysoptionStore.Model, error) {
		mu.RLock()
		defer mu.RUnlock()

		cloned := make([]sysoptionStore.Model, len(items))
		copy(cloned, items)
		return cloned, nil
	})

	if err := svc.StartAutoRefresh(ctx); err != nil {
		t.Fatalf("start auto refresh: %v", err)
	}

	initial, err := svc.GetString(context.Background(), "notice")
	if err != nil {
		t.Fatalf("get initial value: %v", err)
	}
	if initial != "old" {
		t.Fatalf("expected initial value old, got %s", initial)
	}

	mu.Lock()
	items = []sysoptionStore.Model{
		{OptionKey: "notice", OptionValue: "new", Status: sysoptionStore.StatusOnline},
	}
	mu.Unlock()

	deadline := time.Now().Add(400 * time.Millisecond)
	for time.Now().Before(deadline) {
		current, currentErr := svc.GetString(context.Background(), "notice")
		if currentErr == nil && current == "new" {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatal("expected auto refresh to load updated value")
}

func TestNewOptionDoesNotPreloadConfig(t *testing.T) {
	ResetOptionForTest()
	config.ResetForTest()
	appEnv.ResetForTest()

	_ = newOption(time.Minute)

	flags := pflag.NewFlagSet("config", pflag.ContinueOnError)
	config.BindFlags(flags)
	if err := flags.Parse([]string{"--port", "4200"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}
	if err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}
	if got := config.Get().Server.Port; got != 4200 {
		t.Fatalf("expected port 4200 after option service construction, got %d", got)
	}
}

func TestCreateOptionRefreshesCacheImmediately(t *testing.T) {
	ResetOptionForTest()
	db.ResetForTest()
	config.ResetForTest()
	appEnv.ResetForTest()

	t.Setenv("APP_DATABASE_SQLITE_PATH", filepath.Join(t.TempDir(), "option-create.db"))
	if err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}

	_ = db.Get()

	item, err := CreateOption(context.Background(), CreateOptionRequest{
		Key:         "site_name",
		Value:       "Gin Template",
		Description: "站点名称",
		IsPublic:    true,
		Type:        sysoptionStore.TypeString,
		Status:      sysoptionStore.StatusOnline,
	})
	if err != nil {
		t.Fatalf("create option: %v", err)
	}
	if item.OptionKey != "site_name" {
		t.Fatalf("expected created key site_name, got %q", item.OptionKey)
	}
	if item.ID <= 0 {
		t.Fatalf("expected created id > 0, got %d", item.ID)
	}

	value, err := GetOptionString(context.Background(), "site_name")
	if err != nil {
		t.Fatalf("get option after create: %v", err)
	}
	if value != "Gin Template" {
		t.Fatalf("expected refreshed cache value, got %q", value)
	}
}

func TestUpdateOptionRejectsInvalidJSON(t *testing.T) {
	ResetOptionForTest()
	db.ResetForTest()
	config.ResetForTest()
	appEnv.ResetForTest()

	t.Setenv("APP_DATABASE_SQLITE_PATH", filepath.Join(t.TempDir(), "option-invalid-json.db"))
	if err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}

	_ = db.Get()

	_, err := UpdateOption(context.Background(), "about", UpdateOptionRequest{
		Value:       "{invalid",
		Description: "关于信息",
		IsPublic:    true,
		Type:        sysoptionStore.TypeJSON,
		Status:      sysoptionStore.StatusOnline,
	})
	if !errors.Is(err, ErrInvalidOptionJSON) {
		t.Fatalf("expected ErrInvalidOptionJSON, got %v", err)
	}
}

func TestGetOptionStringReturnsNotFoundForOfflineOption(t *testing.T) {
	ResetOptionForTest()
	db.ResetForTest()
	config.ResetForTest()
	appEnv.ResetForTest()

	t.Setenv("APP_DATABASE_SQLITE_PATH", filepath.Join(t.TempDir(), "option-offline.db"))
	if err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}

	_ = db.Get()

	if err := sysoptionStore.Create(context.Background(), &sysoptionStore.Model{
		OptionKey:   "offline_notice",
		OptionValue: "offline",
		Description: "下线配置",
		IsPublic:    true,
		Type:        sysoptionStore.TypeString,
		Status:      sysoptionStore.StatusOffline,
	}); err != nil {
		t.Fatalf("create offline option: %v", err)
	}

	_, err := GetOptionString(context.Background(), "offline_notice")
	if !errors.Is(err, ErrOptionNotFound) {
		t.Fatalf("expected ErrOptionNotFound for offline option, got %v", err)
	}
}
