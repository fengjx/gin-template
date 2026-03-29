---
name: flow-review
description: review the change against spec, tests, and regressions, and produce a findings-first review log.
---

适用场景：

- 实现完成后，需要进行代码审查和规格对照

前置输入：

- `40-implementation.md`
- `50-test-report.md` 可为空但需说明
- `feature.yaml.current_gate=test-review`

必须读取：

- `AGENTS.md`
- `docs/code_review.md`
- `docs/features/<NNN>-<slug>/20-tech-spec.md`
- `docs/features/<NNN>-<slug>/30-tasks.md`
- `docs/features/<NNN>-<slug>/40-implementation.md`
- `docs/features/<NNN>-<slug>/50-test-report.md`

必须写入：

- `docs/features/<NNN>-<slug>/60-review.md`

输出规则：

- findings-first
- 按严重度排序
- 必须指出行为风险、缺失测试、规格漂移、上线风险
- 严重度和输出口径遵循 `docs/code_review.md`
- 如果没有发现问题，也要写残余风险和观察点

退出条件：

- `60-review.md` 已给出是否允许进入发布文档阶段

阻塞规则：

- 发现阻塞问题时，必须进入 `flow-investigate`

下一步建议：

- review 通过且测试通过时进入 `flow-doc-release`
- 任一阻塞问题进入 `flow-investigate`
