# 20 Tech Spec

## 目标与非目标

- 目标：抽取 `pkg/cache` 最小通用缓存组件，并用它重构 `internal/service/option.go`。
- 非目标：扩展 TTL、单 key loader、OpenAPI、数据库 schema、对外配置接口。

## 现状与约束

- `option` 当前自己维护 `map + RWMutex + 定时刷新 goroutine`，并把加载逻辑直接绑定到服务里。
- 仓库要求新增通用能力优先放到 `pkg`，同时业务服务仍需保留清晰的错误语义和测试覆盖。
- goroutine 需要通过 `errs.Recover()` 保护，后台循环应继续沿用 `pkg/timex.SetInterval`。

## 模块边界

- `pkg/cache`：负责内存快照、前缀命名空间、全量刷新和可选后台刷新。
- `internal/service/option.go`：负责系统配置的业务 API、写库逻辑、`GetAndCache` 语义和日志。
- `docs/features/006-generic-cache-refactor/*`：负责记录本次重构的实现、测试与 review 信息。

## 接口与数据流

- `pkg/cache` 对外提供：
  - `Loader[K ~string, V any]`
  - `Options[K ~string, V any]`
  - `New`
  - `Run`
  - `Get`
  - `GetAll`
  - `Set`
  - `Refresh`
  - `Del`
  - `ErrLoaderNotConfigured`
- `Cache` 内部存储使用带前缀的 key；对外 API 和 `GetAll` 始终暴露业务 key。
- `Refresh` 从 loader 获取完整 map 后：
  - 未配置前缀：直接整体替换缓存。
  - 配置前缀：只保留命中此前缀的项，再整体替换缓存。
- `option.GetString` 数据流：
  - 先 `cache.Get(key)`
  - miss 时 `Refresh`
  - 再次 `cache.Get(key)`
  - 若仍 miss，则返回 `errOptionNotFound`

## OpenAPI / 配置 / 数据库影响

- 无 OpenAPI 变更。
- 无新增配置项。
- 无数据库 schema 变更。

## 风险、回滚与观测点

- 风险：前缀过滤或去前缀逻辑不一致，会让对外 key 视图错误。
- 风险：`Run` 的幂等处理不当，会导致重复启动后台循环。
- 回滚：恢复 `pkg/cache`、`internal/service/option.go` 及相关测试到主线版本即可。
- 观测点：`option cache refreshed` 日志、后台刷新失败日志、`pkg/cache` 与 `option` 单测结果、仓库门禁命令结果。

## 测试与验证策略

- 针对 `pkg/cache` 覆盖：基本读写删除、`GetAll` 副本语义、全量刷新替换、前缀过滤、未配置 loader、`Run` 首次刷新/自动刷新/幂等/停止。
- 针对 `option` 覆盖：首次懒加载、写库后即时刷新、后台自动刷新、不提前触发配置加载、缺失 key 错误映射。
- 执行 `make gen`、`make verify`、`make lint`、`make test`、`make check`。
