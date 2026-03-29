package command

import (
	"context"
	"fmt"

	"gin-template/internal/app/config"
	"gin-template/internal/app/db"
	"gin-template/pkg/errs"
	"github.com/spf13/cobra"
)

func newSchemaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Schema 相关命令",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "verify",
		Short: "校验当前数据库 schema",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := config.Load(); err != nil {
				return errs.Wrap(err, "加载配置失败")
			}
			if err := db.VerifySchema(context.Background()); err != nil {
				return errs.Wrap(err, "校验数据库 schema 失败")
			}
			fmt.Fprintln(cmd.OutOrStdout(), "schema verified")
			return nil
		},
	})
	return cmd
}
