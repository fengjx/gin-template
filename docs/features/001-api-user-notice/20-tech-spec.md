# 20 Tech Spec

## 目标与非目标

- 目标：新增只读公告接口
- 非目标：公告管理后台

## 现状与约束

- API 契约必须先行
- 前端只消费 envelope 结构

## 模块边界

- `api/openapi/openapi.yaml`
- `internal/biz/option`
- `admin/src/api/client.ts`

## 接口与数据流

首页请求公告接口 -> 后端读取 option -> 返回 envelope。

## OpenAPI / 配置 / 数据库影响

- OpenAPI 新增公告读取路径
- 不新增配置和表

## 风险、回滚与观测点

- 风险：前端空状态处理遗漏
- 回滚：下线新接口并恢复前端兜底文案

## 测试与验证策略

- `make gen`
- `make verify`
- `make check`
