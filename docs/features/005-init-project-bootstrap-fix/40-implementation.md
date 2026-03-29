# 40 Implementation

## 本轮完成

- 修复 `README.md` 中模板初始化脚本的 raw 地址，改为 `fengjx/gin-template`，并补充模板仓库地址说明。
- 为 README 模板初始化区块补齐 `<!-- template:init:end -->`，确保初始化脚本可以删除模板专属说明。
- 调整 `scripts/init-project.sh`，保留 `.git` 历史，将克隆出的 `origin` 重命名为 `template`，并把 GitHub 仓库地址统一归一化为 HTTPS。
- 更新初始化完成提示，引导用户后续通过 `git fetch template` 同步模板修复。
- 刷新 `admin/package-lock.json`，将 `brace-expansion` 和 `picomatch` 的间接依赖升级到安全版本。

## 修改范围

- `README.md`
- `scripts/init-project.sh`
- `admin/package-lock.json`
- `docs/features/005-init-project-bootstrap-fix/*`

## 关键设计选择

- 保留模板 git 历史，不再删除 `.git`，这样新项目可以直接基于 `template` 远程同步上游修复。
- GitHub 仓库地址只做协议归一化，不自动推断和创建用户自己的 `origin`，避免脚本替用户做高风险猜测。
- 审计问题通过刷新模板锁文件解决，不在初始化脚本里执行临时 `npm audit fix`，保证下游项目首次安装即可落到安全解析结果。
- 非 GitHub 地址保持原样，仅对 `github.com` 的 SSH/HTTP 地址转成 HTTPS，减少对用户自定义仓库的意外改写。

## 新增或更新的测试

- 新增初始化脚本语法检查：`bash -n scripts/init-project.sh`
- 新增基于 GitHub SSH 地址的无安装回归，验证 `.git` 保留与 `template` 远程 HTTPS 归一化。
- 新增基于当前工作树快照模板仓库的无安装回归，验证 README 模板区块会在生成项目中被删除。
- 新增基于当前工作树快照模板仓库的完整初始化回归，验证 `npm install` 后 `npm audit` 为 0 漏洞。

## 未完成项

- 分支已提交并推送，但自动创建 PR 阶段被外部凭证阻塞，需修复 `GH_TOKEN` 或在已登录 GitHub 的浏览器中手动完成建 PR。

## 交给 review 的关注点

- 确认 `normalize_repo_url` 仅覆盖 GitHub 地址的行为符合预期。
- 确认 README 中模板初始化说明与脚本最终行为完全一致。
- 确认锁文件升级未引入前端构建或测试回归。
