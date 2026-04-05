---
name: flow-bugfix
description: run the repo's bugfix workflow end to end when the user explicitly asks to use flow-bugfix.
---

适用场景：

- 用户明确要求使用 `flow-bugfix`
- 需要按仓库 bugfix 快路径推进修复
- 需要从调查、feature 初始化和 gate 恢复开始编排整个修复迭代

触发边界：

- 仅在用户点名 `flow-bugfix` 时使用
- 用户未点名时，回到默认模式，由 AI 自主决定工作方式，但仍必须遵守 `AGENTS.md` 与 `docs/playbooks/feature-lifecycle.md`

必须读取：

- `AGENTS.md`
- `docs/playbooks/feature-lifecycle.md`
- `docs/features/README.md`
- `docs/features/_template/feature.yaml`

入口检查：

1. 先确认本地 `main` 已同步到远端最新状态。
2. 任何写操作前，必须已从最新 `main` 创建或切换到本次迭代专用分支。
3. 若当前仍在 `main` 或共享脏工作树中，先切换到独立分支或 worktree，再继续。

初始化 / 恢复规则：

- 若尚无 `docs/features/<NNN>-<slug>/`，先基于模板创建目录，并填写 `feature.yaml`：
  - `type=bugfix`
  - `status=in_progress`
  - `current_gate=investigate`
  - `branch`、`owner`、`created_at`、`updated_at` 与当前迭代一致
- 初始化 bugfix 目录时，保留模板固定文件，并在 `00-brief.md`、`10-prd.md` 标注本项走 bugfix 快路径，跳过 `flow-idea` / `flow-prd`
- 若 feature 目录已存在，必须从 `feature.yaml.current_gate` 恢复。
- 若当前 gate 对应的上游文档缺失，回退到最早缺失的合法 gate，不允许跳步。

bugfix 推进顺序：

`investigate -> spec -> tasks -> impl -> test-review -> doc-release -> pr -> done`

阶段编排要求：

- `investigate`：进入 `flow-investigate`
- `spec`：进入 `flow-spec`
- `tasks`：进入 `flow-tasks`
- `impl`：进入 `flow-impl`
- `test-review`：先执行 `flow-test`，再执行 `flow-review`
- 若测试或 review 出现阻塞，重新进入 `flow-investigate`，再按调查结论回到 `spec` 或 `impl`
- `doc-release`：进入 `flow-doc-release`
- `pr`：停止自动推进，先向用户确认“当前工作已完成，可以提交 PR”，确认后再进入 `flow-pr`

退出条件：

- bugfix 已进入 `pr` 等待用户确认，或已由 `flow-pr` 推进到 `done`
- 若问题暂缓处理，`feature.yaml.current_gate=stopped`
- `feature.yaml` 与各阶段文档保持一致
