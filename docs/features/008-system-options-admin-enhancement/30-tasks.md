# 30 Tasks

## 任务列表

### T1

- 目标：扩展 `sys_options` 契约、数据库结构与 schema version。
- 修改范围：`api/openapi/openapi.yaml`、`database/bootstrap/*`、`database/upgrade/*`、配置默认值
- 前置依赖：无
- 产出物：`type` / `status` 字段、`POST /options` 契约、schema version 3
- 验证命令：`make gen`、`make verify`

### T2

- 目标：实现系统配置创建、更新校验与在线读取语义。
- 修改范围：`internal/store/sysoption`、`internal/service/option.go`、`internal/biz/option/*`、错误码与集成测试
- 前置依赖：T1
- 产出物：创建接口、JSON 校验、上下线读取控制、后端测试
- 验证命令：`make test`

### T3

- 目标：实现后台新增/编辑与简单 JSON 编辑器。
- 修改范围：`admin/src/pages/OptionsPage.tsx`、`admin/src/components/shared/json-editor.tsx`、`admin/src/api/*`、前端测试
- 前置依赖：T1
- 产出物：新增配置交互、类型切换、JSON 编辑器、前端测试
- 验证命令：`make check`
