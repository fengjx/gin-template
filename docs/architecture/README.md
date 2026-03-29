# 架构速览

这份文档只保留 AI 助手和新加入开发者最常用的上下文。

feature 级过程文档不放在这里，统一放到 `docs/features/<NNN>-<slug>/`；这里仅保留稳定的架构约束。

## 启动链路

1. `cmd/server/main.go` 通过 `internal/app/command` 启动 CLI
2. `serve` 命令加载配置、初始化数据库、确保系统选项和默认管理员存在
3. `internal/app/http.NewEngine` 装配中间件、OpenAPI 文档、pprof、业务路由与前端静态资源
4. `internal/biz/*` 模块按 `api.go`、`model.go`、`service.go` 组织，并通过 `init` 在 `api.go` 中向路由注册中心注入 handler

## 配置加载优先级

配置统一由 `internal/app/config` 管理，优先级从低到高如下：

1. 代码默认值
2. `configs/config.yaml`
3. `configs/config.{env}.yaml`
4. `configs/config.local.yaml`
5. `APP_` 环境变量
6. CLI 参数

典型示例：

- `APP_SERVER_PORT=4000`
- `APP_AUTH_JWT_SECRET=change-me`
- `go run ./cmd/server serve --env dev --port 4000`

## OpenAPI 生成链

- 契约源：`api/openapi/openapi.yaml`
- Go 服务端生成物：`internal/app/http/openapi.gen.go`
- 前端 TS 类型生成物：`admin/src/api/generated.ts`
- 所有接口返回的时间字段统一使用 Unix 秒级时间戳（`int64`），不要返回 RFC3339 / `date-time` 字符串

统一入口：

```bash
make gen
make verify
```

`make verify` 会在生成后检查工作区是否有契约漂移。

## 前后端联调

- 后端默认：`http://localhost:3000`
- 前端默认：`http://localhost:5173`
- Vite 代理转发：
  - `/api`
  - `/docs`
  - `/openapi`

前端请求统一走 `admin/src/api/client.ts`，不要在页面组件里直接调用裸 `fetch`。

## 提交前推荐命令

```bash
make setup
make check
```

如果只改了契约或生成物相关内容，请至少执行：

```bash
make gen
make verify
```
