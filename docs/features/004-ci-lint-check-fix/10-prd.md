# 10 PRD

## 功能目标

- 恢复 GitHub CI 中 `backend-lint` 与 `frontend-check` 的稳定通过。

## 用户角色

- 仓库维护者
- 提交 PR 的开发者

## 典型用户流程

- 开发者提交代码。
- GitHub Actions 执行后端 lint 与前端检查。
- 两个任务均通过，开发者可以继续 review 与合并。

## 范围内 / 范围外

- 范围内：修复导致当前 CI 失败的后端与前端代码问题，补齐必要验证与文档。
- 范围外：放宽 CI 标准、重构无关模块、引入新的 lint 约定。

## 状态与边界条件

- 后端校验以 `golangci-lint-action@v6` + `v1.64.8` 为准。
- 前端校验以 `npm run lint`、`npm run typecheck`、`npm run test` 为准。
- 本轮修复不改变业务接口和用户可见功能语义。

## 异常与失败路径

- 若某项修复后仍有 lint/test 阻塞，则继续记录到 `70-bug-log.md` 并收敛根因。
- 若本地结果与 CI 行为不一致，以 [ci.yml](/Users/fengjianxin/workspaces/my-opensource-project/gin-template/.github/workflows/ci.yml) 为真实执行来源。

## 验收标准

- 后端 lint 复现命令通过。
- 前端 `lint`、`typecheck`、`test` 通过。
- 触达文件的补充验证记录写入 `50-test-report.md`。
- `60-review.md` 无阻塞 finding，或阻塞项已关闭。
