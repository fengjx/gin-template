---
name: flow-pr
description: 在 feature 已完成实现、测试和 review 后，整理相关改动、生成规范 commit message，完成 commit、push，并按仓库 GitHub PR 模板创建 PR。
---

适用场景：

- 需要为当前 feature 收尾并发起 PR
- 需要根据已完成改动整流 commit message、推送分支并填写 PR 模板

前置输入：

- 当前 feature 目录已存在，且 `50-test-report.md`、`60-review.md`、`80-release-doc.md` 至少有基础内容
- 工作树中的待提交改动已经明确范围；若存在无关改动，先剔除或单独说明

必须读取：

- `AGENTS.md`
- `.github/pull_request_template.md`
- `docs/playbooks/feature-lifecycle.md`
- 当前 feature 的 `20-tech-spec.md`
- 当前 feature 的 `50-test-report.md`
- 当前 feature 的 `60-review.md`
- 当前 feature 的 `80-release-doc.md`
- `git status --short`

执行步骤：

1. 确认当前不在 `main`，并且分支与 feature 对应。
2. 按 feature 范围检查工作树，只提交与本次 feature 直接相关的文件。
3. 根据改动内容生成规范 commit message：
   - 优先使用 `feat:` / `fix:` / `refactor:` / `docs:` 等前缀
   - 主题行直接概括用户可感知结果或主要工程收益
   - 若本次同时包含流程文档或 skill，同一提交主题应覆盖主改动，不堆砌文件名
4. `git add` 相关文件并执行 `git commit`。
5. 将当前分支 `git push` 到远端同名分支。
6. 按 `.github/pull_request_template.md` 组织 PR 内容：
   - `Summary` 写清问题、结果和 feature 目录
   - `Checks Run` 勾选实际执行过的命令
   - `Review` 写明是否按 `docs/code_review.md` 检查、结论和残余风险
   - `Docs / Release` 对照 `80-release-doc.md`
   - `GitHub / CI Impact` 明确是否影响 CI、发布链路或外部集成
7. 创建 PR，并把 PR 链接或编号回填到 `80-release-doc.md`。

提交与 PR 约束：

- 不要把与当前 feature 无关的本地改动混进提交。
- 若自动创建 PR 受 `GH_TOKEN`、网络或权限阻塞，必须把阻塞原因写入 `80-release-doc.md`，不能直接宣告完成。
- 若最新代码与已记录的测试结果不一致，先补跑验证，再提交。

产出要求：

- 一个规范 commit
- 远端已存在对应分支
- 一个按模板填写的 GitHub PR
- `80-release-doc.md` 已记录 PR 编号或阻塞原因
