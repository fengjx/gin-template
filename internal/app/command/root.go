package command

import (
	"fmt"
	"os"

	"gin-template/internal/app/config"
	appEnv "gin-template/internal/app/env"
	"gin-template/pkg/errs"
	"github.com/spf13/cobra"
)

var rootCmd = newRootCommand()

func newRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "gin-template",
		Short:            "Gin + React 同构脚手架",
		TraverseChildren: true,
	}

	config.BindFlags(rootCmd.PersistentFlags())
	appEnv.BindFlags(rootCmd.PersistentFlags())
	rootCmd.AddCommand(newConfigCommand())
	rootCmd.AddCommand(newServeCommand())
	rootCmd.AddCommand(newSchemaCommand())
	rootCmd.AddCommand(newSeedAdminCommand())
	rootCmd.AddCommand(newOpenAPICommand())
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "打印版本",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := config.Load(); err != nil {
				return errs.Wrap(err, "加载配置失败")
			}
			fmt.Fprintln(cmd.OutOrStdout(), config.Get().App.Version)
			return nil
		},
	})
	return rootCmd
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
