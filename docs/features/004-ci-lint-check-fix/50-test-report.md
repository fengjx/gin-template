# 50 Test Report

## 触达范围

- `pkg/errs`
- `internal/app/env`
- `internal/app/http`
- `internal/app/command`
- `internal/middleware`
- `admin/src/api`
- `admin/src/utils`

## 执行命令

- `env GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOLANGCI_LINT_CACHE=$(pwd)/.cache/go/golangci-lint ./bin/golangci-lint run`
- `make verify`
- `make check`

## 结果

- 以上命令均通过。
- `make check` 已覆盖 `config verify`、`openapi validate`、`make lint`、`make test`、前端 `typecheck` 与 `build`。

## 未覆盖项

- 未直接从 GitHub Actions 拉取原始 run 日志；本轮依据 [ci.yml](/Users/fengjianxin/workspaces/my-opensource-project/gin-template/.github/workflows/ci.yml) 在本地等价复现并验证。

## 结论

- `通过`
