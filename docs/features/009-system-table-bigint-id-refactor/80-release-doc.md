# 80 Release Doc

## 文档同步项

- README：无需更新
- docs/architecture：已纳入 `database.md` 规范补充
- docs/playbooks：无需更新
- `.env.example`：无需更新
- 示例或文档页：无需更新

## 交付摘要

- 系统表主键统一为自增整数，并同步更新文件/配置对外 `id` 类型。
- 新增 SQLite / MySQL 升级 SQL，将 schema version 提升到 `4`。
- PR：[#12](https://github.com/fengjx/gin-template/pull/12)

## 合并前确认

- [x] feature.yaml.status 已更新为 `done`
- [x] feature.yaml.current_gate 已更新为 `done`
