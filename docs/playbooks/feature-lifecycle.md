# 剧本：Feature 闭环流程

这份剧本面向人和 AI，定义本仓库 feature 的标准生命周期。

`docs/features/<NNN>-<slug>/` 是本仓库唯一的 execution plan 载体，不再额外维护独立 `PLANS.md`。

## 流程

标准 feature：

`flow-idea -> flow-prd -> flow-spec -> flow-tasks -> flow-impl -> flow-test + flow-review -> flow-investigate(如有问题) -> flow-doc-release`

bugfix 快路径：

`flow-investigate -> flow-spec -> flow-tasks -> flow-impl -> flow-test -> flow-review -> flow-doc-release`

## 门禁

1. 先创建 `docs/features/<NNN>-<slug>/`
2. 更新 `feature.yaml.current_gate`
3. 当前门禁未通过前，不进入下一 skill
4. `flow-test` 或 `flow-review` 失败时，必须先写 `70-bug-log.md`
5. `flow-doc-release` 未完成前，`feature.yaml.status` 不能是 `done`

## 文件职责

- `00-brief.md`：问题定义、价值判断、可行性
- `10-prd.md`：用户与功能边界
- `20-tech-spec.md`：技术方案、模块、接口、数据流
- `30-tasks.md`：任务拆解、依赖、验证命令
- `40-implementation.md`：实现记录与未完成项
- `50-test-report.md`：测试矩阵、执行命令、结果
- `60-review.md`：审查意见与风险
- `70-bug-log.md`：问题复现、根因、修复收敛
- `80-release-doc.md`：文档与发布同步项

## 命令矩阵

- 触达 `api/openapi`、生成类型、接口契约：`make gen`、`make verify`、`make check`
- 触达 `internal/`、`cmd/`、`database/`、`internal/store/`：`make test`
- 涉及 schema：`go run ./cmd/server schema verify --env dev`
- 触达 `admin/src/`：`cd admin && npm run lint && npm run test && npm run build`
- 文档 only：跳过代码测试，但仍执行 `flow-review` 与 `flow-doc-release`

## dry-run 场景

- `docs/features/001-api-user-notice`：新增 API 场景
- `docs/features/002-admin-usage-insight`：新增管理页场景
- `docs/features/003-auth-cookie-bugfix`：bugfix 场景

## Skills 入口

- 仓库共享 skills 统一放在 `.agents/skills/`
- `flow-review` 必须同时遵循 `docs/code_review.md`

## 本地 GitHub 操作约定

- 需要由本地 AI 或脚本创建 PR 时，优先从仓库根目录 `.env.local` 读取 `GH_TOKEN`
- `GH_TOKEN` 必须使用短期、最小权限的 fine-grained token，并限制到当前仓库
- `.env.local` 仅用于本地机器，不得提交到仓库
