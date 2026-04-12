# 数据库设计规范

- 必须包含主键字段，非特殊情况都适用自增 id bigint
- 系统表统一使用 `sys_` 前缀
- 字段命名使用下划线
- 普通索引命名使用 idx_ + 字段名首字符，如 idx_a_b，除非遇到命名冲突，则使用全字段名
- 唯一索引使用 uk_ + 字段名首字符，如 uk_a_b，除非遇到命名冲突，则使用全字段名
- 所有表必须包含 `ctime` 和 `utime`
```sql
`utime` datetime not null default now() on update now() comment '更新时间'
`ctime` datetime not null default now() comment '创建时间'
```
- 用户 ID 统一使用字段名 `uid`，类型为 `bigint`
- 新表初始化放到 `database/bootstrap/*.sql`
- 版本升级 SQL 放到 `database/upgrade/{driver}`，不要把升级逻辑写入运行时代码
