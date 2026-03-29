# 00 Brief

## 问题重述

- GitHub CI 中 `backend-lint` 与 `frontend-check` 失败，阻塞主分支与后续 PR 的合并信心。

## 用户与场景

- 维护者需要主分支恢复绿色 CI。
- 贡献者需要在本地与 CI 上得到一致的校验结果。

## 价值判断

- 优先级高。该问题直接影响交付节奏与回归判断。

## 已知事实 / 推断 / 假设

- 已依据 [ci.yml](/Users/fengjianxin/workspaces/my-opensource-project/gin-template/.github/workflows/ci.yml) 在本地复现。
- 后端失败集中在 Go lint 规则收紧后的旧代码与测试文件。
- 前端失败集中在 `Biome` 的格式与导入顺序漂移。
- 当前不需要调整 OpenAPI、配置或数据库。

## 主要风险

- 为了快速消除告警而放宽 lint 规则，会掩盖真实问题。
- 修复测试文件时如果顺手改语义，可能引入非预期行为变化。

## 最小落地路径

- 保持 CI 配置不变，最小化修复后端 lint 问题与前端格式漂移。
- 按 bugfix 快路径补齐调查、实现、验证、review 与 release 记录。

## 结论

- `继续`
