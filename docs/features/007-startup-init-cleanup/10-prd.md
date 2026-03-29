# 10 PRD

## 用户与场景

- 面向维护者：需要在不改变对外行为的前提下，进一步简化服务启动入口。

## 范围

- 收敛 `serve` 启动阶段的 bootstrap 入口。
- 在 `internal/service` 暴露统一 `Init` 入口承接后台初始化。
- 为 `internal/app/log` 增补与现有包装风格一致的 `Panic` / `PanicCtx`。

## 非目标

- 不调整 OpenAPI、数据库 schema、配置文件。
- 不重写整个启动流程，也不引入新的后台任务。
