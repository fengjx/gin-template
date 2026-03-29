# 80 Release Doc

## 文档同步项

- README：无需改动
- docs/architecture：无需改动
- docs/playbooks：无需改动
- `.env.example`：无需改动
- 示例或文档页：无需改动

## 交付摘要

- 修复 GitHub CI 中 `backend-lint` 与 `frontend-check` 的阻塞问题。
- 后端修复聚焦静态检查兼容与测试文件规范化，前端修复聚焦 `Biome` 格式漂移。
- 已补齐本地复现、验证、review 与 release 闭环文档。

## 合并前确认

- [x] feature.yaml.status 已更新为 `done`
- [x] feature.yaml.current_gate 已更新为 `done`
