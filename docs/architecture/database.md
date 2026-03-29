# 数据库设计规范

- 系统表统一使用 `sys_` 前缀
- 所有表必须包含 `ctime` 和 `utime`
- 字段命名使用下划线
- 所有表必须包含 utime 和 ctime
```sql
`utime` datetime not null default now() on update now() comment '更新时间'
`ctime` datetime not null default now() comment '创建时间'
```
- 用户 ID 统一使用字段名 `uid`，类型为 `bigint`
- 新表初始化放到 `database/bootstrap/*.sql`
- 版本升级 SQL 放到 `database/upgrade/{driver}`，不要把升级逻辑写入运行时代码
