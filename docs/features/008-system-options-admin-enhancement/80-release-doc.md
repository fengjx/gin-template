# 80 Release Doc

## 本次变更

- 系统配置后台支持新增配置。
- 管理员可在后台新增和编辑系统配置。
- `sys_options` 新增 `type` 与 `status` 字段。
- JSON 类型配置在后台支持简单 JSON 编辑器。
- 下线配置不会再被 `/system/about`、`/system/notice`、`pprof_url` 等业务读口返回。

## 发布注意事项

- 旧库需要先执行 `database/upgrade/{driver}/20260411_add_sys_options_type_status.sql`。
- PR：[#11](https://github.com/fengjx/gin-template/pull/11)
