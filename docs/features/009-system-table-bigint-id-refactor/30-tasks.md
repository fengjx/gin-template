# 30 Tasks

## 任务列表

### T1

- 目标：调整 bootstrap SQL 和升级 SQL，统一 6 张系统表的主键与 `sys_options` 唯一约束。
- 修改范围：`database/bootstrap/*`、`database/upgrade/*`、`internal/app/config/config.go`
- 前置依赖：无
- 产出物：schema version `4`、新旧库一致的系统表结构、`uk_o`
- 验证命令：`go run ./cmd/server schema verify --env dev`

### T2

- 目标：同步更新 store、业务处理和 OpenAPI/前端类型，确保 `Option.id`、`File.id` 与文件路由参数改为数值型。
- 修改范围：`internal/store/*`、`internal/biz/file/*`、`internal/biz/option/*`、`internal/service/option.go`、`api/openapi/openapi.yaml`、生成产物、`admin/src/api/client.ts`
- 前置依赖：T1
- 产出物：统一的 `int64`/`number` 类型链路
- 验证命令：`make gen`、`make verify`、`make check`

### T3

- 目标：补充实现记录、测试报告和发布文档，完成本地验证与提交。
- 修改范围：`docs/features/009-system-table-bigint-id-refactor/*`
- 前置依赖：T1、T2
- 产出物：实现记录、测试结果、本地提交
- 验证命令：`make lint`、`make test`
