# 80 Release Doc

## 文档同步项

- README：已同步初始化脚本地址、模板仓库地址、`template` 远程说明和模板区块边界。
- docs/architecture：无需变更；本次未修改长期架构约束。
- docs/playbooks：无需变更；现有 feature 生命周期与交付要求未变化。
- `.env.example`：无需变更；本次未新增配置项。
- 示例或文档页：无需变更；本次未涉及文档页或截图资源。

## 交付摘要

- 修复 README 初始化入口错误，新增模板仓库说明。
- 初始化脚本改为保留 git 历史并将模板远程统一为 `template` / HTTPS。
- 刷新 `admin/package-lock.json`，使新项目首次安装后的 `npm audit` 归零。
- 待补充 PR 链接。

## 合并前确认

- [ ] feature.yaml.status 已更新为 `done`
- [ ] feature.yaml.current_gate 已更新为 `done`
