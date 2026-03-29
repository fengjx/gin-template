# 70 Bug Log

## 记录

### BUG-1

- 症状：`backend-lint` 在 GitHub CI 中失败，本地复现可见 `errcheck`、`revive`、`gocritic`、`gosimple` 与 `staticcheck` 报错。
- 复现步骤：执行 `env GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOLANGCI_LINT_CACHE=$(pwd)/.cache/go/golangci-lint ./bin/golangci-lint run`。
- 影响范围：阻塞后端 lint 通过，并影响整体 CI 绿色状态。
- 假设与排查过程：先读取 [ci.yml](/Users/fengjianxin/workspaces/my-opensource-project/gin-template/.github/workflows/ci.yml) 确认真实入口，再本地复现，定位到 `pkg/errs` 等历史代码与测试文件不满足当前 lint 规则，以及 `pflag` 废弃字段兼容问题。
- 根因：仓库现有后端代码未完全适配当前 `golangci-lint` 规则集合与 Go 版本检查。
- 修复建议：在不改变行为的前提下补齐错误处理、收敛无用参数与命名问题、替换废弃 API、清理测试中的旧写法。
- 回归验证点：后端 lint 命令通过，相关测试通过。

### BUG-3

- 症状：PR #6 的 GitHub Actions 中，`backend-lint` 仍失败，日志报错 `the Go language version (go1.24) used to build golangci-lint is lower than the targeted Go version (1.25.0)`。
- 复现步骤：查看 [job 69045922048](https://github.com/fengjx/gin-template/actions/runs/23701506426/job/69045922048?pr=6) 日志。
- 影响范围：PR 无法合并，且本地 `make setup` 安装的 lint 版本也与 CI 漂移。
- 假设与排查过程：先通过 GitHub Actions annotations 与 job logs 确认失败不是新的 lint finding，而是 `golangci-lint-action@v6` 下载的 `v1.64.8` 二进制由 Go 1.24 构建；随后本地安装 `golangci-lint v2.6.2` 并迁移 `.golangci.yml` 验证通过。
- 根因：CI 与本地 setup 都固定在不支持目标 Go 1.25 的 `golangci-lint v1.64.8`。
- 修复建议：将 CI 和本地安装入口同时升级到支持 Go 1.25 的 `golangci-lint v2.6.2`，并将配置文件迁移到 v2 格式。
- 回归验证点：新的 `./bin/golangci-lint run ./...`、`make lint`、`make check` 通过，PR 中 `backend-lint` 重新变绿。

### BUG-4

- 症状：升级到 `golangci-lint v2.6.2` 后，PR 的 `backend-lint` 仍失败，报错 `golangci-lint v2 is not supported by golangci-lint-action v6, you must update to golangci-lint-action v7`；同时 `generate-check` 失败，CI 中生成的 `internal/app/http/openapi.gen.go` 与仓库提交不一致。
- 复现步骤：查看 `backend-lint` 的 [job 69046940143](https://github.com/fengjx/gin-template/actions/runs/23701800575/job/69046940143) 和 `generate-check` 的 [job 69046940149](https://github.com/fengjx/gin-template/actions/runs/23701800575/job/69046940149) 日志。
- 影响范围：PR #6 无法通过所有必需检查。
- 假设与排查过程：先读取 check-run annotations 确认 `backend-lint` 失败不是 lint finding，而是 action major 版本不支持 v2；再读取 `generate-check` 日志，确认 `go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen` 在 CI 中解析到了 `v2.6.0`，与仓库生成文件基于 `v2.5.0` 的结果不一致。
- 根因：工具链升级后只改了 `golangci-lint` 版本，没有同步升级 `golangci-lint-action` major 版本；`openapi generate` 命令没有显式固定 `oapi-codegen` 版本。
- 修复建议：将 `golangci/golangci-lint-action` 升级到 `@v7`，并在 `openapi generate` 中显式使用 `github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.0`。
- 回归验证点：`make verify` 本地无生成漂移，PR 中 `backend-lint` 与 `generate-check` 重新通过。

### BUG-2

- 症状：`frontend-check` 在 GitHub CI 中失败，本地复现为 `Biome` 对 import 顺序与格式输出不一致。
- 复现步骤：执行 `cd admin && npm run lint`。
- 影响范围：阻塞前端检查通过，并影响整体 CI 绿色状态。
- 假设与排查过程：先按 CI 中的 `lint -> typecheck -> test` 顺序本地执行，确认目前 `typecheck` 已通过，主要阻塞集中在 `admin/src/api/client.ts`、`admin/src/api/client.test.ts`、`admin/src/utils/format.ts` 的格式漂移。
- 根因：前端文件提交时未与当前 `Biome` 规则保持一致。
- 修复建议：仅做格式与导入顺序修正，并继续执行 `typecheck`、`test` 验证是否存在次级问题。
- 回归验证点：前端 `lint`、`typecheck`、`test` 通过，必要时补 `build` 验证。
