---
name: flow-tasks
description: break the technical design into execution-ready tasks with scope, dependencies, outputs, and verification commands.
---

适用场景：

- 技术方案已确定，需要转成可执行任务

前置输入：

- `20-tech-spec.md`
- `feature.yaml.current_gate=tasks`

必须读取：

- `AGENTS.md`
- `docs/playbooks/feature-lifecycle.md`
- 相关原子剧本
- `docs/features/<NNN>-<slug>/20-tech-spec.md`

必须写入：

- `docs/features/<NNN>-<slug>/30-tasks.md`

每个任务必须包含：

- 编号
- 标题
- 目标
- 修改范围
- 前置依赖
- 产出物
- 验证命令

方案设计必须符合以下规范：

- 整体架构方案 `docs/architecture/README.md`
- api 设计规范 `docs/architecture/api.md`
- 服务端开发规范 `docs/architecture/backend.md`
- 管理后台开发规范 `docs/architecture/admin.md`
- 数据库设计规范 `docs/architecture/database.md`


退出条件：

- `30-tasks.md` 可以直接交给实现者逐项执行
- `feature.yaml.current_gate` 已更新为 `impl`

阻塞规则：

- 如果任务粒度不可 review、验证命令缺失或依赖顺序不清，不得进入 `flow-impl`

下一步建议：

- 进入 `flow-impl`
