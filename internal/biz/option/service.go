package option

import (
	"errors"
	"fmt"
	"strings"

	"gin-template/internal/app/config"
	appEnv "gin-template/internal/app/env"
	appService "gin-template/internal/service"
	"gin-template/pkg/errs"
	"github.com/gin-gonic/gin"
)

var errPprofDisabled = errors.New("pprof 未启用")
var errPprofURLNotConfigured = errors.New("pprof_url 未配置")

func getOptionValue(c *gin.Context, key string) (string, error) {
	return appService.GetOptionString(c.Request.Context(), key)
}

func buildPprofURL(c *gin.Context, cfg config.Config) (string, error) {
	if !cfg.Server.PprofEnabled {
		return "", errs.WithStack(errPprofDisabled)
	}
	if appEnv.IsDev() {
		return fmt.Sprintf("http://localhost:%d/debug/pprof/", cfg.Server.Port), nil
	}

	value, err := appService.GetOptionString(c.Request.Context(), "pprof_url")
	if err != nil {
		return "", errs.Wrap(err, "读取 pprof_url 配置失败")
	}
	if strings.TrimSpace(value) == "" {
		return "", errs.WithStack(errPprofURLNotConfigured)
	}
	return strings.TrimSpace(value), nil
}
