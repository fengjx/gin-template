# 剧本：新增表或升级 SQL

适用于 `flow-spec`、`flow-tasks`、`flow-impl` 阶段处理新增系统表、补字段或升级 schema。

## 上游输入

- `docs/features/<NNN>-<slug>/20-tech-spec.md`
- `docs/features/<NNN>-<slug>/30-tasks.md`

## 步骤

1. 设计表名，系统表统一使用 `sys_` 前缀
2. 确保包含 `ctime` 与 `utime`
3. 新库初始化 SQL 放到 `database/bootstrap/{driver}.sql`
4. 已发布版本的升级 SQL 放到 `database/upgrade/{driver}`
5. 在 `internal/store/<tablepackage>` 新增或更新 model/store
6. 在业务模块中接入数据访问
7. 执行 `go run ./cmd/server schema verify --env dev`
8. 补充 store 或业务测试，并记录到 `50-test-report.md`

## 注意事项

- 用户 ID 字段名统一为 `uid`
- 不要把历史升级逻辑塞进运行时代码
- MySQL 与 SQLite 差异要分别处理
