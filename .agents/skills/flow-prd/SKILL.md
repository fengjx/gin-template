---
name: flow-prd
description: turn an approved brief into a feature definition with clear scope, flows, states, and acceptance criteria.
---

适用场景：

- `flow-idea` 已通过，需要形成可评审的功能定义

前置输入：

- `00-brief.md`
- `feature.yaml.current_gate=prd`

必须读取：

- `AGENTS.md`
- `docs/architecture/admin.md`
- `docs/architecture/api.md`
- `docs/playbooks/feature-lifecycle.md`
- `docs/features/<NNN>-<slug>/00-brief.md`

必须写入：

- `docs/features/<NNN>-<slug>/10-prd.md`

输出至少包含：

1. 功能目标
2. 用户角色
3. 典型用户流程
4. 范围内 / 范围外
5. 状态与边界条件
6. 异常与失败路径
7. 验收标准

退出条件：

- `10-prd.md` 已明确产品边界和验收口径
- `feature.yaml.current_gate` 已更新为 `spec`

阻塞规则：

- 如果范围内 / 范围外不清、失败路径缺失或验收标准不可测试，不得进入 `flow-spec`

下一步建议：

- 进入 `flow-spec`
