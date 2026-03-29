# 项目宪法

这是一个 Gin + React + Vite 的前后端同构项目，目标是作为现代化、AI 开发友好的开源脚手架。

## 核心原则

- API 契约以 `api/openapi/openapi.yaml` 为唯一来源，所有接口变更先改契约，再做实现与生成。
- 默认保持 SQLite 可开箱运行，MySQL 作为兼容选项存在，不能让新增流程破坏本地快速启动。
- 业务代码优先按 `internal/biz/<module>/api.go`、`model.go`、`service.go` 组织；跨模块复用逻辑放到 `internal/service`。
- 前端统一通过 `admin/src/api/client.ts` 请求后端，不在页面中直接写裸 `fetch`。
- 所有重要改动都必须有测试与文档落点，AI 参与开发也不能绕过验证和可追溯性。

## 仓库地图

- `api/openapi`：OpenAPI 契约与生成配置
- `cmd/server`：服务入口与 CLI 命令
- `internal/app`：配置、HTTP、日志、数据库、启动流程等基础设施
- `internal/biz`：业务模块与路由注册
- `internal/service`：跨 biz 模块复用逻辑
- `internal/store`：数据模型与存储层
- `admin/src`：管理后台页面、组件、路由与 API 客户端
- `docs/architecture`：长期稳定的工程约束
- `docs/playbooks`：原子工程剧本和闭环流程说明
- `docs/features`：feature 级 execution plan 与迭代产物
- `.agents/skills`：仓库共享 skills 的唯一权威目录

## 规范入口

- 后端规范：`docs/architecture/backend.md`
- 前端规范：`docs/architecture/admin.md`
- 数据库规范：`docs/architecture/database.md`
- API 规范：`docs/architecture/api.md`
- 闭环流程：`docs/playbooks/feature-lifecycle.md`
- review 基线：`docs/code_review.md`
- feature 模板：`docs/features/README.md`

## Run The Project

- 初始化依赖：`make setup`
- 启动后端：`make dev-backend`
- 启动前端：`make dev-admin`
- 本地接近 CI 的完整检查：`make ci-local`

默认地址：

- 后端：`http://localhost:3000`
- 前端：`http://localhost:5173`
- 文档：`http://localhost:3000/docs`

## Build / Test / Lint

- 生成契约与类型：`make gen`
- 校验生成物未漂移：`make verify`
- 后端 + 前端检查：`make check`
- lint：`make lint`
- 测试：`make test`
- 前端专项检查：`make admin-check`
- schema 校验：`go run ./cmd/server schema verify --env dev`
- OpenAPI 校验：`go run ./cmd/server openapi validate`
- 契约 fuzz：`make contract-fuzz`

## 工程约束

- 新增接口：先改 `api/openapi/openapi.yaml`，再改后端 handler，再执行 `make gen`
- 新增配置：先补 `internal/app/config/config.go`，再补 `configs/config.yaml`、`.env.example`、README
- 新增页面：先补路由，再补页面组件，再补 API client 或类型，再补测试
- 改数据库：先补 bootstrap 或 upgrade SQL，再补 store 与业务逻辑，再执行 `make test`
- 生成文件禁止手改；错误返回统一走现有 Problem / envelope 约定
- 需要更细的原子流程时，优先引用 `docs/playbooks/*`，不要把长规则再复制到 feature 文档里

## Execution Plan 约定

- 本仓库不额外维护独立 `PLANS.md`
- `docs/features/<NNN>-<slug>/` 就是 execution plan 与实现闭环的唯一载体
- 任何非微小功能改动都必须创建 feature 目录
- 标准流程：`flow-idea -> flow-prd -> flow-spec -> flow-tasks -> flow-impl -> flow-test + flow-review -> flow-investigate(如有问题) -> flow-doc-release`
- bugfix 快路径：`flow-investigate -> flow-spec -> flow-tasks -> flow-impl -> flow-test -> flow-review -> flow-doc-release`

## Definition Of Done

只有同时满足以下条件才算完成：

- 相关 feature 文档已推进到正确阶段，必要时 `feature.yaml.status=done`
- 与改动直接相关的测试、lint、typecheck、构建命令已经执行并记录
- 行为与 `docs/features/<NNN>-<slug>/20-tech-spec.md` 一致，没有未解释的规格漂移
- `docs/code_review.md` 要求的 review 已完成，阻塞问题已关闭或显式记录残余风险
- 文档与配置说明已同步到 `80-release-doc.md`

## PR 交付要求

- PR 必须包含对应的 `docs/features/<NNN>-<slug>/`
- PR 描述必须写清执行过的命令、review 结论、是否更新 `80-release-doc.md`
- 如果改动影响 GitHub 工作流、CI、发布链路或外部集成，要在 PR 中单独说明风险
- review 统一遵循 `docs/code_review.md`

## Worktree / Thread 规则

- 长任务或并行任务优先使用独立 git worktree，避免多个线程共享同一工作目录
- 一个 feature 目录对应一个主线程；并行 agent 只允许处理无重叠修改范围
- 当线程变长或上下文开始漂移时，优先压缩总结或新开线程，不要在一个线程里混多个 feature
- 调试或 review 线程可以 fork，但回写结论时必须回到主 feature 目录

## GitHub / CI 外部上下文

- 当问题涉及 PR、Issue、CI 失败、分支基线、合并风险时，优先读取 GitHub 与 CI 上下文，而不是只看本地工作区
- 当本地命令与 CI 行为冲突时，以 `.github/workflows/ci.yml` 为真实执行来源，并同步修正文档
- 只在外部上下文能消除人工复制粘贴时才扩展 MCP，不要为“可能有用”而预接大量工具

## 提交前最低检查

- `make gen`
- `make verify`
- `make lint`
- `make test`
- `make check`
