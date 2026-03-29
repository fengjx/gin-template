# 80 Release Doc

## 文档同步项

- README：无
- docs/architecture：已更新 `docs/architecture/backend.md`，补充 `internal/service` struct 不使用 `Service` 后缀约束
- docs/playbooks：已更新 `docs/playbooks/feature-lifecycle.md`，明确必须先从最新 `main` 切出专用分支再开始写入
- `.env.example`：无
- 示例或文档页：新增 `docs/features/006-generic-cache-refactor/*`；新增 `.agents/skills/flow-pr/SKILL.md`

## 交付摘要

- 新增 `pkg/cache` 最小通用缓存组件，并重构 `option` 服务缓存实现。
- 同步更新后端命名规范、feature 生命周期约束，并新增 `flow-pr` skill 规范提交与 PR 收尾流程。
- PR：待创建

## 合并前确认

- [ ] feature.yaml.status 已更新为 `done`
- [ ] feature.yaml.current_gate 已更新为 `done`
