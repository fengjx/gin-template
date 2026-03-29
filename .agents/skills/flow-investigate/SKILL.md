---
name: flow-investigate
description: reproduce a failing behavior, form hypotheses, identify the root cause, and log the resolution path before fixing.
---

适用场景：

- 测试失败
- review 提出阻塞问题
- 线上或联调发现 bug

前置输入：

- 具体症状、报错、复现步骤或 findings

必须读取：

- `AGENTS.md`
- `docs/features/<NNN>-<slug>/40-implementation.md`
- `docs/features/<NNN>-<slug>/50-test-report.md`
- `docs/features/<NNN>-<slug>/60-review.md`

必须写入：

- `docs/features/<NNN>-<slug>/70-bug-log.md`

输出至少包含：

1. 症状
2. 复现步骤
3. 影响范围
4. 假设与排查过程
5. 根因
6. 修复建议
7. 回归验证点

退出条件：

- 已有可执行的修复方向
- `feature.yaml.current_gate` 已更新为 `impl`

阻塞规则：

- 未复现或未收敛到可信根因前，不得直接修复

下一步建议：

- 回到 `flow-impl`
