# 30 Tasks

## 任务列表

### T1

- 目标：修复后端 `backend-lint` 报错并保持现有行为不变
- 修改范围：`pkg/errs`、`internal/app/env`、`internal/app/http`、`internal/app/command`、`internal/middleware`
- 前置依赖：`70-bug-log.md` 已记录本地复现与根因
- 产出物：后端代码/测试最小修复
- 验证命令：`env GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOLANGCI_LINT_CACHE=$(pwd)/.cache/go/golangci-lint ./bin/golangci-lint run`

### T2

- 目标：修复前端 `frontend-check` 的格式与潜在测试阻塞
- 修改范围：`admin/src/api`、`admin/src/utils`
- 前置依赖：`70-bug-log.md` 已记录本地复现与根因
- 产出物：前端最小修复
- 验证命令：`cd admin && npm run lint && npm run typecheck && npm run test`

### T3

- 目标：执行回归矩阵并更新闭环文档
- 修改范围：`docs/features/004-ci-lint-check-fix`
- 前置依赖：T1、T2 完成
- 产出物：测试报告、review 结论、release 记录
- 验证命令：`make test`、`cd admin && npm run build`
