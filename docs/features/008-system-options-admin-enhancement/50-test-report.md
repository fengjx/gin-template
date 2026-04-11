# 50 Test Report

## 执行命令

- `make gen`
- `make test`
- `cd admin && npm run test`
- `make lint`
- `make check`
- `make verify`

## 结果

- `make gen`：通过
- `make test`：通过
- `cd admin && npm run test`：通过
- `make lint`：通过
- `make check`：通过（已在最终权限口径与 UI 修正后重新执行）
- `make verify`：已执行；生成步骤通过，但命令按仓库约定会对 `internal/app/http/openapi.gen.go` 和 `admin/src/api/generated.ts` 的未提交改动做 diff，因此在当前 feature 实现阶段返回非零。这两处 diff 均为本次 OpenAPI 契约变更带来的预期生成物更新。

## 覆盖点

- 后端 store：`type` / `status` 字段读写、公开配置过滤
- 后端 service：创建刷新、JSON 非法值拒绝、下线配置读取
- OpenAPI 集成：列表返回扩展字段、创建成功、重复 key、非法 JSON、非法类型、下线后对外读口不可用
- 前端组件：JSON 编辑器错误提示与格式化
- 前端页面：新增配置、JSON 非法时阻断提交
