# 剧本：新增配置项

适用于 `flow-spec`、`flow-tasks`、`flow-impl` 阶段处理新增后端或前端运行配置。

## 上游输入

- `docs/features/<NNN>-<slug>/20-tech-spec.md`
- `docs/features/<NNN>-<slug>/30-tasks.md`

## 步骤

1. 在 `internal/app/config/config.go` 增加结构体字段、默认值和必要归一化逻辑
2. 如果需要 CLI 覆盖，在 `BindFlags` 中绑定参数
3. 在 `configs/config.yaml` 增加基线配置
4. 仅把环境差异项放进 `configs/config.dev.yaml` 或 `configs/config.prod.yaml`
5. 如果是敏感信息，补到 `.env.example`
6. 更新 README 与 `80-release-doc.md` 中的配置说明
7. 补充配置加载优先级测试

## 注意事项

- 保持优先级顺序稳定
- 默认值要保证项目可在 SQLite 开箱运行
- 敏感值不要提交到仓库
