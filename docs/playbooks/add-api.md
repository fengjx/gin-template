# 剧本：新增接口

适用于 `flow-spec`、`flow-tasks`、`flow-impl` 阶段处理新增或调整后端 API。

## 上游输入

- `docs/features/<NNN>-<slug>/20-tech-spec.md`
- `docs/features/<NNN>-<slug>/30-tasks.md`

## 步骤

1. 修改 `api/openapi/openapi.yaml`
2. 执行 `make gen`
3. 在 `internal/biz/<module>/api.go` 中补充路由与 handler，在 `model.go` 中补请求/响应模型，在 `service.go` 中补可复用业务逻辑
4. 通过 `api.go` 中的 `init` 将路由注册到 `internal/app/registry`
5. 如果接口影响前端，更新 `admin/src/api/client.ts` 或直接消费新生成的类型
6. 补充后端测试，必要时补充前端页面测试
7. 在 `50-test-report.md` 中记录 `make verify` 与 `make check` 的执行结果

## 注意事项

- 生成文件禁止手改
- 错误返回统一走 `Problem`
- 需要登录或管理员权限的接口，统一通过中间件控制
- 所有接口返回的时间字段统一使用 Unix 秒级时间戳（`int64`），OpenAPI 中不要声明为 `string/date-time`
