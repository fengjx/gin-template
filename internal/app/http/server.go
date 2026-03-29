package http

import (
	"bytes"
	"fmt"
	"io/fs"
	"mime"
	stdhttp "net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	projectassets "gin-template"
	"gin-template/internal/app/berr"
	"gin-template/internal/app/config"
	"gin-template/internal/app/docs"
	appEnv "gin-template/internal/app/env"
	"gin-template/internal/app/registry"
	"gin-template/internal/middleware"
	"gin-template/pkg/errs"
)

func NewEngine() (*gin.Engine, error) {
	cfg := config.Get()
	if !appEnv.IsDev() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(middleware.Trace(), middleware.AccessLog(), middleware.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CORSAllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Trace-Id"},
		ExposeHeaders:    []string{"X-Trace-Id"},
		AllowCredentials: true,
	}))
	if len(cfg.Server.TrustedProxies) > 0 {
		if err := engine.SetTrustedProxies(cfg.Server.TrustedProxies); err != nil {
			return nil, errs.Wrap(err, "设置受信任代理失败")
		}
	}

	if cfg.Docs.Enabled {
		engine.GET("/openapi/openapi.yaml", func(c *gin.Context) {
			spec, err := projectassets.ReadOpenAPI()
			if err != nil {
				Abort(c, berr.ErrOpenAPISpecUnavailable.WithError(err))
				return
			}
			c.Data(stdhttp.StatusOK, "application/yaml", spec)
		})
		engine.GET("/docs", func(c *gin.Context) {
			c.Data(stdhttp.StatusOK, "text/html; charset=utf-8", []byte(docs.SwaggerHTML("/openapi/openapi.yaml")))
		})
	}

	if cfg.Server.PprofEnabled {
		registerPprof(engine)
	}

	api := engine.Group("/api/v1")
	for _, registerRoute := range registry.Routes() {
		registerRoute(api)
	}

	registerSPA(engine)
	return engine, nil
}

func Serve() error {
	engine, err := NewEngine()
	if err != nil {
		return errs.Wrap(err, "初始化 HTTP 引擎失败")
	}
	cfg := config.Get()
	if err := engine.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		return errs.Wrap(err, "启动 HTTP 监听失败")
	}
	return nil
}

func registerSPA(engine *gin.Engine) {
	adminFS, err := projectassets.AdminAssets()
	if err != nil {
		return
	}

	engine.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") || strings.HasPrefix(c.Request.URL.Path, "/openapi/") || strings.HasPrefix(c.Request.URL.Path, "/debug/pprof") {
			Abort(c, berr.ErrResourceNotFound)
			return
		}

		filename := strings.TrimPrefix(path.Clean(c.Request.URL.Path), "/")
		if filename == "." || filename == "" {
			filename = "index.html"
		}
		if filename == "index.html" {
			serveEmbeddedFile(c, adminFS, filename)
			return
		}
		if file, err := adminFS.Open(filename); err == nil {
			_ = file.Close()
			c.FileFromFS(filename, stdhttp.FS(adminFS))
			return
		}
		serveEmbeddedFile(c, adminFS, "index.html")
	})
}

func serveEmbeddedFile(c *gin.Context, assetFS fs.FS, filename string) {
	data, err := fs.ReadFile(assetFS, filename)
	if err != nil {
		Abort(c, berr.ErrFrontendAssetUnavailable.WithError(err))
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "text/html; charset=utf-8"
	}
	c.Data(stdhttp.StatusOK, contentType, bytes.Clone(data))
}
