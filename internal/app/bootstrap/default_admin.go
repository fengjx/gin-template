package bootstrap

import (
	"context"

	appLog "gin-template/internal/app/log"
	"gin-template/internal/app/security"
	sysuserStore "gin-template/internal/store/sysuser"
	"gin-template/pkg/errs"
	"go.uber.org/zap"
)

const (
	defaultAdminUsername = "admin"
	defaultAdminEmail    = "admin@example.com"
	defaultAdminPassword = "admin"
)

func EnsureDefaultAdmin(ctx context.Context) error {
	total, err := sysuserStore.Count(ctx)
	if err != nil {
		return errs.Wrap(err, "统计现有用户数量失败")
	}
	if total > 0 {
		return nil
	}

	hash, err := security.HashPassword(defaultAdminPassword)
	if err != nil {
		return errs.Wrap(err, "生成默认管理员密码哈希失败")
	}

	item := sysuserStore.New(defaultAdminUsername, defaultAdminEmail, hash)
	item.Role = sysuserStore.RoleAdmin
	item.EmailVerified = true
	if err := sysuserStore.Create(ctx, item); err != nil {
		return errs.Wrap(err, "创建默认管理员失败")
	}

	appLog.InfoCtx(ctx, "default admin account created",
		zap.Int64("uid", item.UID),
		zap.String("username", item.Username),
	)
	return nil
}
