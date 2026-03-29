# 前端管理后台开发规范

- 页面放在 `admin/src/pages`
- 页面路由统一在 `admin/src/app/App.tsx` 维护
- 共用组件优先放在 `admin/src/components`
- API 请求统一通过 `admin/src/api/client.ts`，不要在页面内直接写裸 `fetch`
- 管理台页面新增后要补充至少一个基础测试
- 前端只消费统一 envelope 响应，不要假设接口直接返回业务对象；成功数据统一从 `data` 字段解包
- 新增或调整错误码时，必须同步更新 `admin/src/api/client.ts`、对应测试以及相关页面的提示文案，避免前后端对同一状态码解释不一致
