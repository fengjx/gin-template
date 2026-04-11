# 20 Tech Spec

## 目标与非目标

- 目标：扩展 `sys_options` 元数据，新增创建接口与后台交互，并让业务读取遵守上下线状态。
- 非目标：引入外部 JSON 编辑器依赖、改变现有配置缓存机制、增加配置版本管理。

## 现状与约束

- `sys_options` 当前只有 `option_key`、`option_value`、`description`、`is_public`。
- 后端仅提供 `GET /options` 和 `PUT /options`。
- 管理后台仅支持编辑已有配置，且值输入统一用 `Textarea`。
- 数据库 schema version 当前为 `2`，本次需提升到 `3` 并补升级 SQL。

## 模块边界

- OpenAPI：`api/openapi/openapi.yaml`
- 数据库：`database/bootstrap/*`、`database/upgrade/*`
- 后端：`internal/store/sysoption`、`internal/service/option.go`、`internal/biz/option/*`
- 前端：`admin/src/pages/OptionsPage.tsx`、`admin/src/components/shared/json-editor.tsx`、`admin/src/api/*`

## 接口定义与数据字段

- `Option` 响应新增：
  - `type: string`，枚举 `string` / `json`
  - `status: string`，枚举 `online` / `offline`
- 新增 `POST /options`
- `PUT /options` 请求体扩展 `type`、`status`
- `sys_options` 新增字段：
  - `type`：默认 `string`
  - `status`：默认 `online`

## 数据流

- 后台列表：直接读取 `sys_options` 全量数据，不做上下线过滤，后台新增和更新统一允许管理员操作。
- 创建/更新：service 校验 `key`、`type`、`status` 和 JSON 合法性，写库后刷新配置缓存。
- 业务读取：仍走 option cache，但仅返回 `status=online` 的配置值。

## 风险、回滚与观测点

- 风险：schema version 升级后，旧库若未执行升级 SQL 将无法通过 schema 校验。
- 回滚：回退 `sys_options` 扩展字段相关代码与 OpenAPI 变更，并将 schema version 退回。
- 观测点：`option cache refreshed` 日志、OpenAPI 集成测试、OptionsPage JSON 编辑器测试、schema verify 结果。

## 测试与验证策略

- `make gen`
- `make verify`
- `make lint`
- `make test`
- `make check`
- 必要时补 `cd admin && npm run test`
