# 后端开发规范

- 所有代码都需要有详细的注释说明
- 所有模块优先通过 Go 原生 `init` 组织模块注册与初始化顺序
- model 按表拆分 package，去掉下划线并统一全小写，例如 `sys_user` 使用 `sysuser`
- 全局异常必须通过中间件记录详细错误栈
- `applog` 使用 `log.app.filename` 控制输出目标：配置文件名时固定输出到文件且使用 JSON；未配置时固定输出到控制台且使用空格分隔格式。`accesslog` 使用 `log.access.filename` 指定输出文件且固定使用 JSON。两类文件日志都必须按自然天切割，并通过 `trace_id` 关联
- 配置优先级固定为：默认值 < `configs/config.yaml` < `configs/config.{env}.yaml` < `configs/config.local.yaml` < 环境变量 < CLI
- 环境变量统一使用 `APP_` 前缀，例如 `APP_SERVER_PORT`、`APP_AUTH_JWT_SECRET`
- `internal/store/<table>` 存放数据库访问相关代码，所有方法需要有完整的单元测试
- `internal/biz` 的业务逻辑直接写到当前目录下，可服用的代码单独一个文件编写，必须包含完整测单元测试
- `internal/biz/<module>` 固定拆分为 `api.go`、`model.go`、`service.go`：`api.go` 负责路由注册与 handler，`model.go` 负责请求/响应模型与转换，`service.go` 负责可复用业务逻辑
- `internal/service` 为可复用的逻辑封装，解决各 biz 模块之间的功能依赖，可服用的逻辑都放到 `internal/service` 并以自己的 package 来命名文件，例如 option 相关逻辑复用，则创建 `internal/service/option.go` 文件来编写，直接使用普通函数封装即可，需要有完整的单元测试
- 接口参数和返回值使用下划线命名规范
- 所有 error 处理，都使用 errs 进行包装，业务异常统一在 berr 中定义，并由 `http.Abort` 统一处理
- 新增错误时，先确定归属是系统级还是某个 biz；若属于 biz，必须落到对应号段，避免跨模块复用同一业务码
- 需要附带底层错误时，统一使用 `WithError`；需要覆盖面向用户的明细时，统一使用 `WithDetail`
- 所有 goroutine 执行的方法都需要使用 defer errs.Recover() 来保护
- 数据库常量字段，在代码中必须要定义

## 错误码规范

- 所有业务错误模板统一定义在 `internal/app/berr/errors.go`，禁止在 biz handler 或 middleware 中临时拼装新的业务码
- `internal/app/berr/errors.go` 中的错误码和 `BusinessError` 模板必须按区域组织：先定义系统级状态码，再按 `auth`、`user`、`file`、`option` 等 biz 模块分区
