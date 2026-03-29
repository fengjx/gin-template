package command

import (
	"os"
	"os/exec"

	appOpenAPI "gin-template/internal/app/openapi"
	"gin-template/pkg/errs"
	"github.com/spf13/cobra"
)

func newOpenAPICommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "openapi",
		Short: "OpenAPI 工具",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "validate",
		Short: "校验 OpenAPI 规范",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := appOpenAPI.ValidateEmbeddedSpec(); err != nil {
				return errs.Wrap(err, "校验 OpenAPI 规范失败")
			}
			_, _ = cmd.OutOrStdout().Write([]byte("openapi validated\n"))
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "generate",
		Short: "生成 OpenAPI 相关代码",
		RunE: func(cmd *cobra.Command, _ []string) error {
			goCmd := exec.Command("go", "run", "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen", "-config", "api/openapi/codegen.yaml", "api/openapi/openapi.yaml")
			goCmd.Stdout = cmd.OutOrStdout()
			goCmd.Stderr = cmd.ErrOrStderr()
			goCmd.Env = withToolEnv()
			if err := goCmd.Run(); err != nil {
				return errs.Wrap(err, "生成 Go OpenAPI 代码失败")
			}

			npmCmd := exec.Command("npm", "run", "generate:api")
			npmCmd.Stdout = cmd.OutOrStdout()
			npmCmd.Stderr = cmd.ErrOrStderr()
			npmCmd.Dir = "admin"
			npmCmd.Env = withToolEnv()
			if err := npmCmd.Run(); err != nil {
				return errs.Wrap(err, "生成前端 API 代码失败")
			}
			return nil
		},
	})
	return cmd
}

func withToolEnv() []string {
	env := os.Environ()
	env = append(env, "GOPROXY=https://proxy.golang.org,direct")
	return env
}
