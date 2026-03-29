# 30 Tasks

## 任务列表

### T1

- 目标：新增页面与路由
- 修改范围：`admin/src/pages`、`admin/src/app/App.tsx`
- 前置依赖：无
- 产出物：页面骨架和路由
- 验证命令：`cd admin && npm run test`

### T2

- 目标：补导航入口和共享组件
- 修改范围：`admin/src/components/AppShell.tsx`
- 前置依赖：T1
- 产出物：菜单入口
- 验证命令：`cd admin && npm run lint`

### T3

- 目标：补前端测试与构建验证
- 修改范围：页面测试
- 前置依赖：T1、T2
- 产出物：测试用例与构建结果
- 验证命令：`cd admin && npm run lint && npm run test && npm run build`
