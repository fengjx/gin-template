package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	appEnv "gin-template/internal/app/env"
	"gin-template/pkg/errs"
)

type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Trace     TraceConfig     `mapstructure:"trace"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Storage   StorageConfig   `mapstructure:"storage"`
	Mail      MailConfig      `mapstructure:"mail"`
	OAuth     OAuthConfig     `mapstructure:"oauth"`
	Turnstile TurnstileConfig `mapstructure:"turnstile"`
	Docs      DocsConfig      `mapstructure:"docs"`
	OpenAPI   OpenAPIConfig   `mapstructure:"openapi"`
	Feature   FeatureConfig   `mapstructure:"feature"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type ServerConfig struct {
	Host             string   `mapstructure:"host"`
	Port             int      `mapstructure:"port"`
	TrustedProxies   []string `mapstructure:"trusted_proxies"`
	FrontendURL      string   `mapstructure:"frontend_url"`
	CORSAllowOrigins []string `mapstructure:"cors_allow_origins"`
	PprofEnabled     bool     `mapstructure:"pprof_enabled"`
}

type DatabaseConfig struct {
	Driver        string `mapstructure:"driver"`
	DSN           string `mapstructure:"dsn"`
	SQLitePath    string `mapstructure:"sqlite_path"`
	SchemaVersion int    `mapstructure:"schema_version"`
}

type TraceConfig struct {
	HeaderName string `mapstructure:"header_name"`
}

type AuthConfig struct {
	Issuer                string `mapstructure:"issuer"`
	AccessTokenTTLMinutes int    `mapstructure:"access_token_ttl_minutes"`
	RefreshTokenTTLHours  int    `mapstructure:"refresh_token_ttl_hours"`
	JWTSecret             string `mapstructure:"jwt_secret"`
	CookieDomain          string `mapstructure:"cookie_domain"`
	SecureCookie          bool   `mapstructure:"secure_cookie"`
	SameSite              string `mapstructure:"same_site"`
}

type StorageConfig struct {
	LocalDir string `mapstructure:"local_dir"`
}

type MailConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	FromName  string `mapstructure:"from_name"`
	FromEmail string `mapstructure:"from_email"`
	SMTPHost  string `mapstructure:"smtp_host"`
	SMTPPort  int    `mapstructure:"smtp_port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
}

type OAuthProviderConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
	AppID        string `mapstructure:"app_id"`
	AppSecret    string `mapstructure:"app_secret"`
}

type OAuthConfig struct {
	GitHub OAuthProviderConfig `mapstructure:"github"`
	WeChat OAuthProviderConfig `mapstructure:"wechat"`
}

type TurnstileConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	SecretKey string `mapstructure:"secret_key"`
}

type DocsConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type OpenAPIConfig struct {
	ValidateOnBoot bool `mapstructure:"validate_on_boot"`
}

type RateLimitWindow struct {
	Requests      int `mapstructure:"requests"`
	WindowSeconds int `mapstructure:"window_seconds"`
}

type RateLimitConfig struct {
	Enabled  bool            `mapstructure:"enabled"`
	Global   RateLimitWindow `mapstructure:"global"`
	Critical RateLimitWindow `mapstructure:"critical"`
	Upload   RateLimitWindow `mapstructure:"upload"`
}

type FeatureConfig map[string]bool

var (
	cfg         Config
	v           = newViper()
	loadOnce    sync.Once
	loadErr     error
	configPath  string
	loadedFiles []string
)

func newViper() *viper.Viper {
	instance := viper.New()
	instance.SetConfigType("yaml")
	setDefaults(instance)
	instance.SetEnvPrefix("APP")
	instance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	instance.AutomaticEnv()
	return instance
}

func BindFlags(flags *pflag.FlagSet) {
	flags.StringVar(&configPath, "config", "", "配置文件路径")
	flags.String("host", "", "服务监听地址")
	flags.Int("port", 0, "服务端口")
	flags.String("database-driver", "", "数据库驱动")
	flags.String("database-dsn", "", "数据库 DSN")
	flags.String("sqlite-path", "", "SQLite 数据库路径")
	flags.String("jwt-secret", "", "JWT 密钥")
	flags.String("frontend-url", "", "前端地址")
	flags.Bool("pprof-enabled", false, "是否启用 pprof")

	_ = v.BindPFlag("server.host", flags.Lookup("host"))
	_ = v.BindPFlag("server.port", flags.Lookup("port"))
	_ = v.BindPFlag("database.driver", flags.Lookup("database-driver"))
	_ = v.BindPFlag("database.dsn", flags.Lookup("database-dsn"))
	_ = v.BindPFlag("database.sqlite_path", flags.Lookup("sqlite-path"))
	_ = v.BindPFlag("auth.jwt_secret", flags.Lookup("jwt-secret"))
	_ = v.BindPFlag("server.frontend_url", flags.Lookup("frontend-url"))
	_ = v.BindPFlag("server.pprof_enabled", flags.Lookup("pprof-enabled"))
}

func Load() error {
	loadOnce.Do(func() {
		loadedFiles = nil
		if configPath != "" {
			if err := mergeConfigFile(configPath); err != nil {
				loadErr = errs.Wrap(err, "读取配置文件失败")
				return
			}
		} else {
			for _, candidate := range candidateConfigFiles() {
				if _, err := os.Stat(candidate); err == nil {
					if err := mergeConfigFile(candidate); err != nil {
						loadErr = errs.Wrap(err, "读取配置文件失败")
						return
					}
				}
			}
		}

		if err := v.Unmarshal(&cfg); err != nil {
			loadErr = errs.Wrap(err, "解析配置失败")
			return
		}
		normalize()
		if err := Validate(cfg); err != nil {
			loadErr = err
		}
	})
	return loadErr
}

func MustLoad() {
	if err := Load(); err != nil {
		panic(err)
	}
}

func Get() Config {
	MustLoad()
	return cfg
}

func Sources() []string {
	MustLoad()
	return slices.Clone(loadedFiles)
}

func Validate(current Config) error {
	if strings.TrimSpace(current.App.Name) == "" {
		return errs.WithStack(errors.New("config invalid: app.name 不能为空"))
	}
	if current.Server.Port <= 0 || current.Server.Port > 65535 {
		return errs.WithStack(errors.New("config invalid: server.port 必须在 1-65535 之间"))
	}
	switch strings.ToLower(current.Auth.SameSite) {
	case "default", "lax", "strict", "none":
	default:
		return errs.WithStack(errors.New("config invalid: auth.same_site 仅支持 default、lax、strict、none"))
	}
	if strings.EqualFold(current.Database.Driver, "mysql") && strings.TrimSpace(current.Database.DSN) == "" {
		return errs.WithStack(errors.New("config invalid: database.driver=mysql 时必须提供 database.dsn"))
	}
	if strings.TrimSpace(current.Database.SQLitePath) == "" {
		return errs.WithStack(errors.New("config invalid: database.sqlite_path 不能为空"))
	}
	if len(current.Server.CORSAllowOrigins) == 0 {
		return errs.WithStack(errors.New("config invalid: server.cors_allow_origins 至少配置一个来源"))
	}
	return nil
}

func ResetForTest() {
	configPath = ""
	cfg = Config{}
	loadOnce = sync.Once{}
	loadErr = nil
	loadedFiles = nil
	v = newViper()
}

func setDefaults(target *viper.Viper) {
	target.SetDefault("app.name", "gin-template")
	target.SetDefault("app.version", "0.1.0")
	target.SetDefault("server.host", "0.0.0.0")
	target.SetDefault("server.port", 3000)
	target.SetDefault("server.trusted_proxies", []string{})
	target.SetDefault("server.frontend_url", "http://localhost:5173")
	target.SetDefault("server.cors_allow_origins", []string{"http://localhost:5173"})
	target.SetDefault("server.pprof_enabled", true)
	target.SetDefault("database.driver", "sqlite")
	target.SetDefault("database.dsn", "")
	target.SetDefault("database.sqlite_path", "runtime/data/app.db")
	target.SetDefault("database.schema_version", 3)
	target.SetDefault("trace.header_name", "X-Trace-Id")
	target.SetDefault("auth.issuer", "gin-template")
	target.SetDefault("auth.access_token_ttl_minutes", 30)
	target.SetDefault("auth.refresh_token_ttl_hours", 720)
	target.SetDefault("auth.jwt_secret", "change-me")
	target.SetDefault("auth.cookie_domain", "")
	target.SetDefault("auth.secure_cookie", false)
	target.SetDefault("auth.same_site", "lax")
	target.SetDefault("storage.local_dir", "runtime/uploads")
	target.SetDefault("docs.enabled", true)
	target.SetDefault("openapi.validate_on_boot", true)
	target.SetDefault("feature", map[string]bool{})
	target.SetDefault("rate_limit.enabled", true)
	target.SetDefault("rate_limit.global.requests", 120)
	target.SetDefault("rate_limit.global.window_seconds", 60)
	target.SetDefault("rate_limit.critical.requests", 10)
	target.SetDefault("rate_limit.critical.window_seconds", 300)
	target.SetDefault("rate_limit.upload.requests", 30)
	target.SetDefault("rate_limit.upload.window_seconds", 300)
}

func candidateConfigFiles() []string {
	files := []string{"configs/config.yaml"}
	currentEnv := appEnv.Current()
	if currentEnv != "" {
		files = append(files, fmt.Sprintf("configs/config.%s.yaml", currentEnv))
	}
	files = append(files, "configs/config.local.yaml")
	return files
}

func normalize() {
	if cfg.Database.DSN == "" {
		cfg.Database.Driver = "sqlite"
	}
	if cfg.Database.Driver == "" {
		cfg.Database.Driver = "sqlite"
	}
	if cfg.Server.FrontendURL == "" && len(cfg.Server.CORSAllowOrigins) > 0 {
		cfg.Server.FrontendURL = cfg.Server.CORSAllowOrigins[0]
	}
	if len(cfg.Server.CORSAllowOrigins) == 0 && cfg.Server.FrontendURL != "" {
		cfg.Server.CORSAllowOrigins = []string{cfg.Server.FrontendURL}
	}
	cfg.Database.SQLitePath = filepath.Clean(cfg.Database.SQLitePath)
	cfg.Storage.LocalDir = filepath.Clean(cfg.Storage.LocalDir)
}

func mergeConfigFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errs.Wrap(err, "读取配置文件内容失败")
	}
	if err := v.MergeConfig(bytes.NewReader(data)); err != nil {
		return errs.Wrap(err, "合并配置文件失败")
	}
	loadedFiles = append(loadedFiles, path)
	return nil
}
