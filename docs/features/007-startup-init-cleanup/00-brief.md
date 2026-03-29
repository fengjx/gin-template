# 00 Brief

## 背景

- `serve` 命令启动阶段同时承担配置加载、bootstrap 和后台服务初始化，入口职责偏散。
- `internal/service` 已经有 `option` 单例，但外部仍直接调用 `StartOptionAutoRefresh`，不利于后续收敛更多服务级初始化。

## 问题

- 启动链路缺少统一的 service init 入口，后续若继续扩展初始化流程，`serve.go` 会持续膨胀。
- `internal/app/log` 缺少包级 `Panic` 包装，启动期错误日志无法和现有 `Info/Error` 风格保持一致。

## 结论

- 作为一次小范围启动整理，保留现有行为不变，只收敛入口与日志 API，便于后续继续扩展初始化阶段。
