---
name: flow-test
description: execute the required verification matrix based on the touched areas and write a concrete test report.
---

适用场景：

- 实现完成后，需要执行功能测试、回归测试与命令矩阵

前置输入：

- `40-implementation.md`
- `feature.yaml.current_gate=test-review`

必须读取：

- `AGENTS.md`
- `docs/playbooks/feature-lifecycle.md`
- `docs/features/<NNN>-<slug>/20-tech-spec.md`
- `docs/features/<NNN>-<slug>/30-tasks.md`
- `docs/features/<NNN>-<slug>/40-implementation.md`

必须写入：

- `docs/features/<NNN>-<slug>/50-test-report.md`

命令矩阵：

- 触达 `api/openapi`、生成类型、接口契约：`make gen`、`make verify`、`make check`
- 触达 `internal/`、`cmd/`、`database/`、`internal/store/`：`make test`
- 涉及 schema：`go run ./cmd/server schema verify --env dev`
- 触达 `admin/src/`：`cd admin && npm run lint && npm run test && npm run build`

退出条件：

- `50-test-report.md` 记录了执行命令、结果、未覆盖项与结论

阻塞规则：

- 不允许只写“理论已验证”
- 任一关键命令失败时，必须把结论标记为阻塞，并进入 `flow-investigate`

下一步建议：

- 测试通过后等待 review 结果
- 测试失败时进入 `flow-investigate`
