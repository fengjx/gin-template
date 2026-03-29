# 20 Tech Spec

## 目标与非目标

- 目标：新增前端聚合页
- 非目标：新增专用分析接口

## 现状与约束

- 复用现有 API client
- 页面不得直接写裸 `fetch`

## 模块边界

- `admin/src/pages`
- `admin/src/app/App.tsx`
- `admin/src/components/AppShell.tsx`

## 接口与数据流

页面加载时并行请求现有统计接口，前端聚合展示。

## OpenAPI / 配置 / 数据库影响

- 无 OpenAPI、配置、数据库改动

## 风险、回滚与观测点

- 风险：多接口并发导致加载态复杂
- 回滚：下线菜单入口

## 测试与验证策略

- `cd admin && npm run lint && npm run test && npm run build`
