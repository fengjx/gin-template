package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"gin-template/internal/app/config"
	"gin-template/internal/app/security"
	sysuserStore "gin-template/internal/store/sysuser"
	"gin-template/pkg/errs"
)

func newSeedAdminCommand() *cobra.Command {
	var username string
	var email string
	var password string
	var role string

	cmd := &cobra.Command{
		Use:   "seed admin",
		Short: "创建管理员账号",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := config.Load(); err != nil {
				return errs.Wrap(err, "加载配置失败")
			}
			if username == "" || email == "" || password == "" {
				return errs.WithStack(errors.New("username、email、password 不能为空"))
			}
			if _, err := sysuserStore.ByEmail(context.Background(), email); err == nil {
				return errs.WithStack(errors.New("邮箱已存在"))
			}
			hash, err := security.HashPassword(password)
			if err != nil {
				return errs.Wrap(err, "生成管理员密码哈希失败")
			}
			item := sysuserStore.New(username, email, hash)
			item.Role = role
			item.EmailVerified = true
			if err := sysuserStore.Create(context.Background(), item); err != nil {
				return errs.Wrap(err, "创建管理员账号失败")
			}
			fmt.Fprintf(cmd.OutOrStdout(), "admin created: %d\n", item.UID)
			return nil
		},
	}

	cmd.Flags().StringVar(&username, "username", "root", "用户名")
	cmd.Flags().StringVar(&email, "email", "root@example.com", "邮箱")
	cmd.Flags().StringVar(&password, "password", "", "密码")
	cmd.Flags().StringVar(&role, "role", sysuserStore.RoleRoot, "角色")
	return cmd
}
