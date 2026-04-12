package bootstrap

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	appLog "gin-template/internal/app/log"
	sysoptionStore "gin-template/internal/store/sysoption"
	"gin-template/pkg/errs"
)

type defaultOption struct {
	Key         string
	Value       string
	Description string
	IsPublic    bool
	Type        string
	Status      string
}

var defaultSystemOptions = []defaultOption{
	{Key: "notice", Value: "欢迎使用 gin-template", Description: "系统公告", IsPublic: true, Type: sysoptionStore.TypeString, Status: sysoptionStore.StatusOnline},
	{Key: "about", Value: "Gin + React 同构脚手架", Description: "关于信息", IsPublic: true, Type: sysoptionStore.TypeString, Status: sysoptionStore.StatusOnline},
	{Key: "pprof_url", Value: "/debug/pprof/", Description: "pprof 监控地址", IsPublic: false, Type: sysoptionStore.TypeString, Status: sysoptionStore.StatusOnline},
}

func EnsureSystemOptions(ctx context.Context) error {
	for _, item := range defaultSystemOptions {
		if _, err := sysoptionStore.ByKey(ctx, item.Key); err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.Wrap(err, "检查系统配置是否已存在失败")
		}

		record := &sysoptionStore.Model{
			OptionKey:   item.Key,
			OptionValue: item.Value,
			Description: item.Description,
			IsPublic:    item.IsPublic,
			Type:        item.Type,
			Status:      item.Status,
		}
		if err := sysoptionStore.Create(ctx, record); err != nil {
			return errs.Wrap(err, "创建系统默认配置失败")
		}
		appLog.InfoCtx(ctx, "system option created", zap.String("key", item.Key))
	}
	return nil
}
