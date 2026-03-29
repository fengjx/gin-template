# api 设计规范

- `api/openapi/openapi.yaml` 是唯一 API 契约源
- `internal/app/http/openapi.gen.go` 与 `admin/src/api/generated.ts` 只能通过生成命令更新，禁止手改
- 新增或修改接口时，必须同步更新 OpenAPI 并执行 `make gen` 与 `make verify`
- 所有接口必须包含完整的集成测试用例，且不依赖具体代码实现，以 api 定义为准
- 当`api/openapi/openapi.yaml`中的定义与代码实现冲突时，需要同步修改
- 所有 JSON 接口统一返回 envelope 结构：成功响应返回 `status`、`msg`、`data`，失败响应额外返回 `details`

## 错误码定义

- `status` 含义固定：`0` 表示成功；`100000-199999` 表示系统级状态码；`200000-299999` 表示业务级状态码
- 业务级状态码必须按 biz 模块划分号段，并在契约、后端实现、前端消费处保持一致；当前约定：
  - `auth` 使用 `200000-209999`
  - `user` 使用 `210000-219999`
  - `file` 使用 `220000-229999`
  - `option` 使用 `230000-239999`
