# Feature 工作区

`docs/features/` 是项目自有 AI 闭环流程的工作区。

每个非微小功能改动都需要创建一个目录：

`docs/features/<NNN>-<slug>/`

例如：

- `docs/features/001-api-user-notice`
- `docs/features/002-admin-usage-insight`
- `docs/features/003-auth-cookie-bugfix`

## 固定文件

每个 feature 目录固定包含：

- `feature.yaml`
- `00-brief.md`
- `10-prd.md`
- `20-tech-spec.md`
- `30-tasks.md`
- `40-implementation.md`
- `50-test-report.md`
- `60-review.md`
- `70-bug-log.md`
- `80-release-doc.md`

## `feature.yaml` 字段

- `id`
- `slug`
- `title`
- `type`：`feature` / `bugfix` / `refactor`
- `status`
- `branch`
- `owner`
- `created_at`
- `updated_at`
- `current_gate`

## 生命周期

标准 feature 流程：

`flow-idea -> flow-prd -> flow-spec -> flow-tasks -> flow-impl -> flow-test + flow-review -> flow-investigate(如有问题) -> flow-doc-release`

bugfix 快路径：

`flow-investigate -> flow-spec -> flow-tasks -> flow-impl -> flow-test -> flow-review -> flow-doc-release`

## 使用方式

- 先阅读 `AGENTS.md`
- 再阅读 `docs/playbooks/feature-lifecycle.md`
- 然后从 `docs/features/_template/` 复制模板开始
- 每一步只写当前 skill 负责的文件，不跳步，不混写
