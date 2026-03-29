# 30 Tasks

## 任务列表

### T1

- 目标：实现 `pkg/cache` 最小通用缓存组件。
- 修改范围：`pkg/cache/cache.go`
- 前置依赖：`20-tech-spec.md`
- 产出物：泛型缓存结构、前缀命名空间、全量刷新、`Run` 生命周期管理、错误定义
- 验证命令：`go test ./pkg/cache -run TestCache`

### T2

- 目标：使用通用缓存重构 `option` 服务内部实现。
- 修改范围：`internal/service/option.go`
- 前置依赖：T1
- 产出物：基于 `pkg/cache` 的 `option`、服务层 `GetAndCache` 语义与现有对外函数兼容
- 验证命令：`go test ./internal/service -run TestOption`

### T3

- 目标：补齐 `pkg/cache` 和 `option` 的测试覆盖。
- 修改范围：`pkg/cache/cache_test.go`、`internal/service/option_test.go`
- 前置依赖：T1-T2
- 产出物：前缀、刷新、自动刷新、懒加载、错误路径相关单测
- 验证命令：`go test ./pkg/cache ./internal/service`

### T4

- 目标：完成本次 refactor 的实现记录与测试记录。
- 修改范围：`docs/features/006-generic-cache-refactor/*`
- 前置依赖：T1-T3
- 产出物：实现记录、测试报告、review 结论、发布说明占位
- 验证命令：`make gen && make verify && make lint && make test && make check`
