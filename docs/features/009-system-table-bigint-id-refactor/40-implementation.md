# 40 Implementation

## 本轮完成

- 调整 6 张系统表的 bootstrap 结构与升级 SQL，统一改为自增整数主键。
- 将 `Option.id`、`File.id` 以及文件接口路径参数同步切换为数值类型。
- 更新 store、业务处理、前端客户端与相关测试。

## 修改范围

- `database/bootstrap/*`
- `database/upgrade/*`
- `internal/store/*`
- `internal/biz/file/*`
- `internal/biz/option/*`
- `internal/service/option.go`
- `api/openapi/openapi.yaml`
- 生成产物与 `admin/src/api/client.ts`

## 关键设计选择

- 升级 SQL 采用直接重建表，显式放弃历史数据迁移。
- `sys_options` 继续以 `option_key` 作为业务主键，`id` 仅作为数据库主键和返回字段。
- `sys_files` 只清数据库记录，不清理历史上传目录。

## 新增或更新的测试

- 更新系统配置 store / service 测试中的 `id` 构造方式。
- 更新前端系统配置页面测试中的 `id` 断言类型。
- 通过命令矩阵覆盖 schema、生成物、后端与前端回归。

## 未完成项

- 待执行完整命令矩阵并写入 `50-test-report.md`。
- 待完成 review、release doc 收尾与 PR。

## 交给 review 的关注点

- 升级 SQL 与 bootstrap 是否完全一致。
- 文件接口的 `id` 解析与错误处理是否满足契约要求。
- `sys_options` 默认数据插入是否仍符合初始化预期。
