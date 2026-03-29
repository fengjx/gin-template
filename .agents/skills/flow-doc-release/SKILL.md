---
name: flow-doc-release
description: finish the loop by syncing docs, examples, and release notes with the shipped behavior.
---

适用场景：

- 实现、测试、review 均已通过，需要完成文档和发布同步

前置输入：

- `50-test-report.md`
- `60-review.md`
- `feature.yaml.current_gate=doc-release`

必须读取：

- `AGENTS.md`
- `README.md`
- `docs/architecture/*`
- `docs/playbooks/*`
- `docs/features/<NNN>-<slug>/20-tech-spec.md`
- `docs/features/<NNN>-<slug>/60-review.md`

必须写入：

- `docs/features/<NNN>-<slug>/80-release-doc.md`

检查清单：

- README 是否需要同步
- 架构约束是否变化
- 原子剧本是否变化
- 配置示例与 `.env.example` 是否变化
- 文档页或示例截图是否需要变化
- PR 是否已自动创建，编号或链接是否已写入 `80-release-doc.md`

退出条件：

- `80-release-doc.md` 已列出所有同步项和完成状态
- PR 已自动创建，且编号或链接已记录到 `80-release-doc.md`
- `feature.yaml.status=done`
- `feature.yaml.current_gate=done`

阻塞规则：

- 只要存在未同步文档项，就不得结束 feature
- 若自动创建 PR 因 `GH_TOKEN`、权限或网络受阻，必须先记录阻塞原因，且不得结束 feature

下一步建议：

- 进入 PR review、合并或发布准备
