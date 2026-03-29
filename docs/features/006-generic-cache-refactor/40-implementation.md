# 40 Implementation

## 本轮完成

- 新增 `pkg/cache`，提供 `New`、`Run`、`Get`、`GetAll`、`Set`、`Refresh`、`Del` 和 `ErrLoaderNotConfigured`。
- 将 `internal/service/option.go` 改为围绕通用缓存做业务封装，移除原来的缓存锁、加载状态和手写刷新循环。
- 保留 `option` 对外 API，并在服务层补回 miss 时刷新一次的读取策略。
- 按最新后端命名规范去掉 `option` 内部结构和测试中的 `Service` 后缀，并同步测试基建重置入口。
- 新增 `flow-pr` skill，用于规范提交、推送和按模板发起 PR 的收尾流程。
- 新增并更新 `pkg/cache` 与 `option` 的单测，并完成仓库最低检查命令回归。

## 修改范围

- `pkg/cache/cache.go`
- `pkg/cache/cache_test.go`
- `internal/service/option.go`
- `internal/service/option_test.go`
- `t/testkit/harness.go`
- `.agents/skills/flow-pr/SKILL.md`
- `docs/features/006-generic-cache-refactor/*`

## 关键设计选择

- 通用层只抽取“全量 loader 驱动的内存快照”能力，不引入业务特定的 `GetAndCache`。
- 前缀在缓存内部形成命名空间，对外 API 始终暴露业务 key，避免调用方重复处理前缀。
- `option` 的刷新日志下沉到 loader 适配层，保持自动刷新与手动刷新都能复用同一套观测点。
- PR 收尾流程抽成独立 skill，避免每次人工重复拼 commit/PR 模板内容。

## 新增或更新的测试

- 新增 `pkg/cache` 单测，覆盖基本读写、前缀过滤、副本语义、未配置 loader、自动刷新、幂等和停止行为。
- 更新 `option` 单测，验证首次读取懒加载、写库后刷新、后台刷新和构造期无副作用。
- 执行 `make gen`、`make verify`、`make lint`、`make test`、`make check`。

## 未完成项

- 暂无。

## 交给 review 的关注点

- 确认 `pkg/cache` 的前缀语义与 `GetAll` 暴露视图一致。
- 确认 `option` 的 miss 刷新策略不会破坏现有错误语义与日志行为。
- 确认最小通用接口没有过早承载 `option` 专属语义，后续其他模块可以直接复用。
