# 剧本：Feature 闭环流程

这份剧本面向人和 AI，定义本仓库 feature 的标准生命周期。

`docs/features/<NNN>-<slug>/` 是本仓库唯一的 execution plan 载体，不再额外维护独立 `PLANS.md`。

## 流程

开始新一轮需求前：

`同步 main 最新状态 -> 创建/切换新分支 -> 创建 docs/features/<NNN>-<slug>/ -> 进入对应 skill 流程`

- 这里的“创建/切换新分支”必须发生在任何实现动作之前，包括创建 feature 文档目录、修改代码、补测试或更新文档；禁止在 `main` 上直接开始执行。

标准 feature：

`flow-idea -> flow-prd -> flow-spec -> flow-tasks -> flow-impl -> flow-test + flow-review -> flow-investigate(如有问题) -> flow-doc-release`

bugfix 快路径：

`flow-investigate -> flow-spec -> flow-tasks -> flow-impl -> flow-test -> flow-review -> flow-doc-release`

## 门禁

1. 新需求开始前，必须先同步本地 `main` 到远端最新状态，再开始新的 feature/bugfix 迭代
2. 必须从最新 `main` 创建或切换到本次迭代的专用分支后，才能创建 `docs/features/<NNN>-<slug>/` 或开始任何写操作
3. 再创建 `docs/features/<NNN>-<slug>/`
4. 更新 `feature.yaml.current_gate`
5. 当前门禁未通过前，不进入下一 skill
6. `flow-test` 或 `flow-review` 失败时，必须先写 `70-bug-log.md`
7. `flow-doc-release` 未完成前，`feature.yaml.status` 不能是 `done`
8. `flow-doc-release` 收尾时必须自动创建 PR，并将 PR 编号或链接写入 `80-release-doc.md`
9. 若自动创建 PR 因 `GH_TOKEN`、权限或网络受阻，必须在 `80-release-doc.md` 记录阻塞原因，且不得结束 feature

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

- 新需求开始前，优先执行一次 `main` 同步，确保新分支基于最新主线
- feature / bugfix 默认使用独立分支承载实现；除紧急只读排查外，不在 `main` 上直接进行需求执行
- 需要由本地 AI 或脚本创建 PR 时，优先从仓库根目录 `.env.local` 读取 `GH_TOKEN`
- `GH_TOKEN` 必须使用短期、最小权限的 fine-grained token，并限制到当前仓库
- `.env.local` 仅用于本地机器，不得提交到仓库
