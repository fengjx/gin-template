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
- `make verify`：通过（已在提交后对生成物一致性再次确认）

## 覆盖点

- 后端 store：`type` / `status` 字段读写、公开配置过滤
- 后端 service：创建刷新、JSON 非法值拒绝、下线配置读取
- OpenAPI 集成：列表返回扩展字段、创建成功、重复 key、非法 JSON、非法类型、下线后对外读口不可用
- 前端组件：JSON 编辑器错误提示与格式化
- 前端页面：新增配置、JSON 非法时阻断提交
