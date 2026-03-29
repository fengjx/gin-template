# 20 Tech Spec

## 目标与非目标

- 目标：通过最小代码修复恢复 `backend-lint` 与 `frontend-check`。
- 非目标：修改 CI 工作流、调整 lint 规则、扩展任何业务需求。

## 现状与约束

- GitHub CI 入口以 [ci.yml](/Users/fengjianxin/workspaces/my-opensource-project/gin-template/.github/workflows/ci.yml) 为准。
- 后端触达 `internal/`、`pkg/` 与 `cmd/` 时，需要补 `make test`。
- 前端触达 `admin/src/` 时，需要补 `cd admin && npm run lint && npm run test && npm run build`。

## 模块边界

- 后端：`pkg/errs`、`internal/app/env`、`internal/app/http`、`internal/app/command`、`internal/middleware`
- 前端：`admin/src/api`、`admin/src/utils`

## 接口与数据流

- 无接口或数据流变更。
- 后端修复聚焦 fmt/io 写入、命名、废弃 API 与测试辅助代码。
- 前端修复聚焦 Biome 规范要求，不改变 API client 语义。

## OpenAPI / 配置 / 数据库影响

- 无 OpenAPI 变更
- 无配置变更
- 无数据库变更

## 风险、回滚与观测点

- 风险：测试辅助代码调整可能误伤断言表达。
- 回滚：逐文件回退本次 lint 修复，并保留 `70-bug-log.md` 的根因与复现记录。
- 观测点：本地 lint/test/build 与 CI 结果保持一致。

## 测试与验证策略

- `env GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOLANGCI_LINT_CACHE=$(pwd)/.cache/go/golangci-lint ./bin/golangci-lint run`
- `make test`
- `cd admin && npm run lint`
- `cd admin && npm run typecheck`
- `cd admin && npm run test`
- `cd admin && npm run build`
