# 20 Tech Spec

## 目标与非目标

- 目标：统一系统表主键规范，提升 schema version 到 `4`，并保持数据库、后端、OpenAPI、前端类型一致。
- 非目标：不保留旧数据，不清理 `runtime/uploads` 历史文件，不新增按 `id` 写配置能力。

## 现状与约束

- `sys_refresh_tokens`、`sys_options`、`sys_oauth_bindings`、`sys_email_verifications`、`sys_password_resets`、`sys_files` 当前都使用字符串主键。
- 项目数据库规范要求主键优先使用自增 `id bigint`。
- `sys_files`、`sys_options` 的 `id` 已暴露到 API 契约与前端类型。

## 模块边界

- 数据库：`database/bootstrap/*`、`database/upgrade/*`
- 后端：对应 `internal/store/*`、`internal/biz/file/*`、`internal/biz/option/*`、`internal/service/option.go`
- 契约与前端：`api/openapi/openapi.yaml`、生成产物、`admin/src/api/client.ts`

## 接口与数据流

- 文件上传后由数据库生成数值型 `id`，返回给后台列表与详情接口。
- 系统配置新增后由数据库生成数值型 `id`，但读写仍以 `option_key` 作为业务定位键。
- 刷新令牌、密码重置、邮箱验证、OAuth 绑定等内部系统表改为数据库分配主键，不再在应用层生成 UUID 主键。

## OpenAPI / 配置 / 数据库影响

- OpenAPI：`Option.id`、`File.id` 以及 `/files/{id}` 参数改为 `integer(int64)`。
- 配置：`database.schema_version` 默认值从 `3` 升到 `4`。
- 数据库：
  - 6 张系统表主键改为自增整数。
  - `sys_options.option_key` 使用命名唯一约束 `uk_o`。
  - 升级 SQL 直接重建目标表并回填默认系统配置。

## 风险、回滚与观测点

- 风险：旧库若未执行升级 SQL，会因为 schema version 不匹配而无法启动。
- 回滚：回退本次代码与 bootstrap/upgrade SQL，并将 schema version 目标值恢复到 `3`。
- 观测点：schema verify、文件增删查、系统配置新增/读取、认证辅助表写入流程。

## 测试与验证策略

- 执行 `make gen`、`make verify`、`make lint`、`make test`、`make check`。
- 执行 `go run ./cmd/server schema verify --env dev` 和 `go run ./cmd/server openapi validate`。
- 使用 SQLite 临时库验证 bootstrap 后 schema version 与表结构。
