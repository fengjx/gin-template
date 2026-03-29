---
name: flow-idea
description: evaluate a raw idea, challenge its framing, and write the feature brief with a clear go/no-go decision.
---

适用场景：

- 新功能想法
- 现有流程优化
- 需要先判断值不值得做的需求讨论

前置输入：

- 原始需求、问题描述或目标
- 目标 feature 目录已创建

必须读取：

- `AGENTS.md`
- `docs/playbooks/feature-lifecycle.md`
- `docs/features/<NNN>-<slug>/feature.yaml`

必须写入：

- `docs/features/<NNN>-<slug>/00-brief.md`

输出至少包含：

1. 问题重述
2. 用户与场景
3. 价值判断
4. 已知事实 / 推断 / 假设
5. 主要风险
6. 最小落地路径
7. 结论：继续 / 暂缓 / 放弃

退出条件：

- `00-brief.md` 已能支撑是否进入 PRD
- `feature.yaml.current_gate` 已更新为 `prd` 或 `stopped`

阻塞规则：

- 如果问题定义不清、目标用户不清或没有最小落地路径，不得进入 `flow-prd`

下一步建议：

- 结论为继续时进入 `flow-prd`
- 结论为暂缓或放弃时只更新 `feature.yaml`，停止后续阶段
