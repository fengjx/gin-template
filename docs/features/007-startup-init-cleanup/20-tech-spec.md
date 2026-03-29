# 20 Tech Spec

## 目标与非目标

- 目标：统一 `serve` 启动中的 bootstrap 与 service init 入口，并补齐包级 panic 日志包装。
- 非目标：改变启动顺序语义、调整系统配置刷新策略、引入新的配置项或数据库变更。

## 现状与约束

- `newServeCommand` 当前直接串联配置加载、bootstrap、配置自动刷新和 HTTP 启动。
- `internal/service` 只有 `option` 单例的具体能力，没有统一初始化入口。
- 启动阶段错误仍需保留当前“立即失败并打印日志”的语义。

## 模块边界

- `internal/app/command/serve.go`：负责启动阶段 orchestration。
- `internal/service/init.go`：负责聚合服务层启动期初始化。
- `internal/app/log/log.go`：负责提供与现有包级日志函数对齐的 panic 包装。

## 接口与数据流

- `serve` 启动流程调整为：
  - `config.Load`
  - 可选 `OpenAPI.ValidateEmbeddedSpec`
  - `db.Get`
  - `boot(ctx)` 负责系统配置和默认管理员 bootstrap
  - `service.Init(ctx)` 负责启动系统配置自动刷新
  - `http.Serve`
- `internal/service` 新增 `Init(ctx)`，当前只封装 `defaultOption.StartAutoRefresh(ctx)`，后续新增初始化逻辑时继续向这里收敛。

## OpenAPI / 配置 / 数据库影响

- 无 OpenAPI 变更。
- 无新增配置项。
- 无数据库 schema 变更。

## 风险、回滚与观测点

- 风险：如果 `Init` 入口后续被滥用，可能重新把业务初始化耦合到服务层聚合入口。
- 回滚：恢复 `serve.go` 对 `option` 自动刷新的直接调用，并删除新增包装函数即可。
- 观测点：`make test` 中 `internal/app/command`、`internal/app/log`、`internal/service` 相关测试结果。

## 测试与验证策略

- `make gen`
- `make verify`
- `make lint`
- `make test`
- `make check`
