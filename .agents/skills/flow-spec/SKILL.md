---
name: flow-spec
description: convert the feature definition into a repo-specific technical design grounded in current architecture and constraints.
---

适用场景：

- PRD 已明确，需要形成工程方案
- bugfix 快路径中的轻量修复方案

前置输入：

- feature：`10-prd.md`
- bugfix：`70-bug-log.md` 中已有复现和根因
- `feature.yaml.current_gate=spec`

必须读取：

- `AGENTS.md`
- `docs/architecture/backend.md`
- `docs/architecture/admin.md`
- `docs/architecture/api.md`
- `docs/architecture/database.md`
- `docs/playbooks/feature-lifecycle.md`
- 相关原子剧本

必须写入：

- `docs/features/<NNN>-<slug>/20-tech-spec.md`

输出至少包含：

- 目标与非目标
- 现状与约束
- 模块边界
- 接口定义与数据字段定义
- 风险、回滚与观测点
- 测试与验证策略

退出条件：

- `20-tech-spec.md` 已明确需要改哪些层、为什么改、怎么验证
- `feature.yaml.current_gate` 已更新为 `tasks`

阻塞规则：

- 如果模块边界不清、验证策略为空或回滚方式缺失，不得进入 `flow-tasks`

下一步建议：

- 进入 `flow-tasks`
