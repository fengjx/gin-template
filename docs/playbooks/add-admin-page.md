# 剧本：新增管理页

适用于 `flow-spec`、`flow-tasks`、`flow-impl` 阶段处理新增后台页面、菜单和 API 联动。

## 上游输入

- `docs/features/<NNN>-<slug>/10-prd.md`
- `docs/features/<NNN>-<slug>/20-tech-spec.md`
- `docs/features/<NNN>-<slug>/30-tasks.md`

## 步骤

1. 在 `admin/src/pages` 新增页面组件
2. 在 `admin/src/app/App.tsx` 注册路由
3. 如果需要导航入口，更新 `admin/src/components/AppShell.tsx`
4. 通过 `admin/src/api/client.ts` 复用统一请求层
5. 优先复用 `admin/src/components/ui` 与 `admin/src/components/shared`
6. 补充页面测试，至少覆盖渲染或关键交互
7. 在 `50-test-report.md` 中记录 `cd admin && npm run lint && npm run test && npm run build`

## 注意事项

- 页面不要直接写裸 `fetch`
- 错误提示统一走现有反馈组件
- 路由保护交给现有认证上下文和受保护布局
