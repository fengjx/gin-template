package command

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"gin-template/internal/app/bootstrap"
	"gin-template/internal/app/config"
	"gin-template/internal/app/db"
	appEnv "gin-template/internal/app/env"
	appHTTP "gin-template/internal/app/http"
	appLog "gin-template/internal/app/log"
	appOpenAPI "gin-template/internal/app/openapi"
	appService "gin-template/internal/service"
	"gin-template/pkg/errs"
)

func newServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "启动服务",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()
			if err := config.Load(); err != nil {
				appLog.Panic("load config error", zap.Error(err))
			}
			if config.Get().OpenAPI.ValidateOnBoot {
				if err := appOpenAPI.ValidateEmbeddedSpec(); err != nil {
					return errs.Wrap(err, "启动前校验 OpenAPI 规范失败")
				}
			}
			_ = db.Get()
			if err := boot(ctx); err != nil {
				appLog.Panic("boot error", zap.Error(err))
			}
			if err := appService.Init(ctx); err != nil {
				appLog.Panic("init service error", zap.Error(err))
			}
			defer appLog.Sync()
			appLog.Info("server 启动", zap.String("env", string(appEnv.Current())))
			if err := appHTTP.Serve(); err != nil {
				return errs.Wrap(err, "启动 HTTP 服务失败")
			}
			return nil
		},
	}
}

func boot(ctx context.Context) error {
	if err := bootstrap.EnsureSystemOptions(ctx); err != nil {
		return errs.Wrap(err, "初始化系统配置失败")
	}
	if err := bootstrap.EnsureDefaultAdmin(ctx); err != nil {
		return errs.Wrap(err, "初始化默认管理员失败")
	}
	return nil
}
