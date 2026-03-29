# 40 Implementation

## 本轮完成

- 已修复后端 `backend-lint` 的静态检查问题，包括 `io.WriteString` 返回值处理、未使用参数、过时 `pflag` API、测试辅助代码的 `gocritic`/`unused` 告警。
- 已修复前端 `frontend-check` 的 `Biome` 格式漂移，包括 import 顺序、换行与长表达式排版。
- 已将 `golangci-lint` 升级到支持 Go 1.25 的 v2.6.2，并把 `.golangci.yml` 迁移到 v2 配置格式，消除 GitHub Actions 上的工具版本不兼容。
- 已修复后续 CI 回归：`golangci-lint-action@v6` 升级到 `@v7`，并将 `openapi generate` 显式固定到 `oapi-codegen@v2.5.0`，避免生成物漂移。
- 已完成 `make verify` 与 `make check` 回归。

## 修改范围

- 后端：`pkg/errs/error.go`、`pkg/errs/stack.go`、`pkg/errs/format_test.go`、`pkg/errs/error_test.go`、`internal/middleware/problem.go`、`internal/app/command/serve.go`、`internal/app/http/resp_test.go`、`internal/app/http/openapi_error_handler_test.go`、`internal/app/env/env.go`
- 前端：`admin/src/api/client.ts`、`admin/src/api/client.test.ts`、`admin/src/utils/format.ts`
- 工具链：`.golangci.yml`、[Makefile](/Users/fengjianxin/workspaces/my-opensource-project/gin-template/Makefile)、[ci.yml](/Users/fengjianxin/workspaces/my-opensource-project/gin-template/.github/workflows/ci.yml)

## 关键设计选择

- 保持 CI 工作流不变，只修代码，不通过放宽规则掩盖问题。
- 对生成接口要求的 `...Id` 命名，仅在测试桩上使用定点 `//nolint:revive`，避免污染其他文件。
- 对仅供测试复用的 helper，使用最小范围 `//nolint:unused` 保留现有断言工具，不改变测试语义。

## 新增或更新的测试

- 更新了 `pkg/errs` 与 `admin/src/api/client.test.ts` 相关测试文件的实现形式，使其符合当前 lint/formatter 规则。
- 未新增测试用例，现有 Go 测试与前端 Vitest 用例已回归通过。

## 未完成项

- 无

## 交给 review 的关注点

- `openapi_error_handler_test.go` 中的 `//nolint:revive` 是否保持在最小范围。
- `pkg/errs/format_test.go` 中保留的 helper 是否仍有后续复用价值。
