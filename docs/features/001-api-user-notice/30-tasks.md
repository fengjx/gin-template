# 30 Tasks

## 任务列表

### T1

- 目标：补 OpenAPI 契约和生成物
- 修改范围：`api/openapi`、生成类型
- 前置依赖：无
- 产出物：接口契约与生成文件
- 验证命令：`make gen && make verify`

### T2

- 目标：实现后端读取逻辑与测试
- 修改范围：`internal/biz/option`
- 前置依赖：T1
- 产出物：handler、service、测试
- 验证命令：`make test`

### T3

- 目标：接入首页公告读取
- 修改范围：`admin/src`
- 前置依赖：T1
- 产出物：首页调用与空状态处理
- 验证命令：`cd admin && npm run lint && npm run test && npm run build`
