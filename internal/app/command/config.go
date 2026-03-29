package command

import (
	"fmt"
	"strings"

	"gin-template/internal/app/config"
	"gin-template/pkg/errs"
	"github.com/spf13/cobra"
)

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "配置相关命令",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "verify",
		Short: "校验当前配置",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := config.Load(); err != nil {
				return errs.Wrap(err, "加载配置失败")
			}
			fmt.Fprintln(cmd.OutOrStdout(), "config verified")
			sources := config.Sources()
			if len(sources) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "config sources: defaults only")
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "config sources: %s\n", strings.Join(sources, ", "))
			return nil
		},
	})
	return cmd
}
