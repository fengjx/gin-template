# 40 Implementation

## 本轮完成

- 扩展 `sys_options` 结构，新增 `type` 与 `status` 字段，并将 schema version 提升到 `3`。
- 更新 OpenAPI 契约，新增 `POST /options`，并扩展 `Option` 与 `PUT /options` 的字段定义。
- 在后端实现系统配置创建能力、`json` 值校验和“仅在线配置可被业务读取”的读取语义。
- 在管理后台新增系统配置创建入口、类型切换、上下线管理和简单 JSON 编辑器。
- 根据最新权限口径，将系统配置的新增和更新权限统一收敛为管理员可操作。
- 补充 store、service、OpenAPI 集成测试、JSON 编辑器测试和 OptionsPage 页面测试。

## 修改范围

- `api/openapi/openapi.yaml`
- `database/bootstrap/*`
- `database/upgrade/*/20260411_add_sys_options_type_status.sql`
- `internal/store/sysoption/*`
- `internal/service/option.go`
- `internal/biz/option/*`
- `internal/app/bootstrap/system_options.go`
- `internal/app/berr/errors.go`
- `admin/src/api/*`
- `admin/src/components/shared/JsonEditor.*`
- `admin/src/pages/OptionsPage.*`
- `docs/features/008-system-options-admin-enhancement/*`

## 关键设计选择

- `type` 使用字符串枚举 `string` / `json`，避免再引入额外的值模式字段。
- 配置缓存继续维护全量快照，但 `GetOptionString` 在读取时过滤 `offline` 状态，保证后台管理和业务消费职责分离。
- JSON 编辑器采用仓库内轻量实现，基于 `Textarea` 增强实时校验和格式化，不引入 Monaco 或 CodeMirror。
- 创建和更新统一走 service 校验 `key`、`type`、`status` 和 JSON 合法性，避免校验逻辑散落在 handler 中。

## 新增或更新的测试

- 新增 `internal/store/sysoption/model_test.go`，覆盖新增字段读写与公开配置过滤。
- 更新 `internal/service/option_test.go`，覆盖创建刷新、非法 JSON 拒绝和下线配置读取。
- 更新 `t/openapi_integration_test.go`，覆盖 `POST /options`、元数据展示、无效输入和下线读口行为。
- 新增 `admin/src/components/shared/JsonEditor.test.tsx`。
- 新增 `admin/src/pages/OptionsPage.test.tsx`。

## 未完成项

- `60-review.md` 尚未写入正式 review 结论。
- `80-release-doc.md` 尚未进入发布收尾与 PR 信息补充阶段。

## 交给 review 的关注点

- 确认 `offline` 语义当前仅影响业务读取，不会误伤后台管理场景。
- 确认 `type=json` 的前后端校验和错误提示足够清晰。
- 确认 schema version 升级与手工 SQL 的发布顺序符合现网运维预期。
