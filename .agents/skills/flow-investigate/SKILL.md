---
name: flow-investigate
description: reproduce a failing behavior, form hypotheses, identify the root cause, and log the resolution path before fixing.
---

适用场景：

- 新 bugfix 进入快路径时的第一步
- 测试失败
- review 提出阻塞问题
- 线上或联调发现 bug

前置输入：

- 具体症状、报错、复现步骤或 findings
- 若为新 bugfix 入口，`feature.yaml.current_gate=investigate`
- 若为实现后返工，需说明失败命令或阻塞 findings

必须读取：

- `AGENTS.md`
- `docs/playbooks/feature-lifecycle.md`
- `docs/features/<NNN>-<slug>/feature.yaml`
- 若来自实现后返工，再读 `docs/features/<NNN>-<slug>/40-implementation.md`
- 若已有测试或 review 结论，再读 `docs/features/<NNN>-<slug>/50-test-report.md`、`docs/features/<NNN>-<slug>/60-review.md`
- 若已有方案或任务，再读 `docs/features/<NNN>-<slug>/20-tech-spec.md`、`docs/features/<NNN>-<slug>/30-tasks.md`

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
8. 建议回流 gate：`spec` 或 `impl`

回流规则：

1. 新 bugfix 入口完成调查后，默认回到 `flow-spec`
2. 测试失败或 review 阻塞，且修复范围局限在既有实现内时，回到 `flow-impl`
3. 若根因改变了方案边界、接口、数据流或任务拆解，回到 `flow-spec`

退出条件：

- 已有可执行的修复方向
- `feature.yaml.current_gate` 已根据调查结论更新为 `spec` 或 `impl`

阻塞规则：

- 未复现或未收敛到可信根因前，不得直接修复

下一步建议：

- 若 `feature.yaml.current_gate=spec`，回到 `flow-spec`
- 若 `feature.yaml.current_gate=impl`，回到 `flow-impl`
