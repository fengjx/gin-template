# 40 Implementation

## 已完成

- 在 `internal/app/command/serve.go` 中提取 `boot(ctx)`，把启动期 bootstrap 收敛到独立函数。
- 在 `internal/service/init.go` 新增 `Init(ctx)` 聚合入口，承接配置自动刷新初始化。
- 在 `internal/app/log/log.go` 新增 `Panic` / `PanicCtx` 包级包装，统一启动阶段错误出口。
- 删除 `internal/service/option.go` 的包级 `StartOptionAutoRefresh` 直接导出，改由 `service.Init` 统一承接。

## 未完成

- 无。
