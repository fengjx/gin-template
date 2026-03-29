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

### BUG-2

- 症状：`frontend-check` 在 GitHub CI 中失败，本地复现为 `Biome` 对 import 顺序与格式输出不一致。
- 复现步骤：执行 `cd admin && npm run lint`。
- 影响范围：阻塞前端检查通过，并影响整体 CI 绿色状态。
- 假设与排查过程：先按 CI 中的 `lint -> typecheck -> test` 顺序本地执行，确认目前 `typecheck` 已通过，主要阻塞集中在 `admin/src/api/client.ts`、`admin/src/api/client.test.ts`、`admin/src/utils/format.ts` 的格式漂移。
- 根因：前端文件提交时未与当前 `Biome` 规则保持一致。
- 修复建议：仅做格式与导入顺序修正，并继续执行 `typecheck`、`test` 验证是否存在次级问题。
- 回归验证点：前端 `lint`、`typecheck`、`test` 通过，必要时补 `build` 验证。
