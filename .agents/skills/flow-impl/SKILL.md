---
name: flow-impl
description: implement one or more approved tasks, keep the change minimal, and update the implementation log.
---

适用场景：

- 已有明确任务拆解，进入代码生产阶段

前置输入：

- `30-tasks.md`
- `feature.yaml.current_gate=impl`

必须读取：

- `AGENTS.md`
- 相关架构规范
- 相关原子剧本
- `docs/features/<NNN>-<slug>/30-tasks.md`

必须写入：

- `docs/features/<NNN>-<slug>/40-implementation.md`

输出至少包含：

1. 本轮完成的任务
2. 修改范围
3. 关键设计选择
4. 新增或更新的测试
5. 未完成项
6. 交给 review 的关注点

退出条件：

- 当前任务已有实现记录
- `feature.yaml.current_gate` 已更新为 `test-review`

阻塞规则：

- 如果任务未闭合、测试未补齐或实现与 `20-tech-spec.md` 明显偏离，必须先修正，不得进入 `flow-test` 或 `flow-review`

下一步建议：

- 并行进入 `flow-test` 与 `flow-review`
