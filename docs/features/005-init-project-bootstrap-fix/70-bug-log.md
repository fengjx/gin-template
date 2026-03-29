# 70 Bug Log

## 记录

### 1. 初始化入口与 git 远程行为异常

- 症状：README 中 raw 地址错误，初始化脚本执行后删除 `.git`，导致用户无法保留模板历史与同步上游修复。
- 复现步骤：阅读 README 初始化命令并检查 `scripts/init-project.sh` 中克隆后的清理逻辑。
- 影响范围：所有依赖 README 初始化项目的新用户。
- 假设与排查过程：先核对 README 中 raw 地址和模板区块，再检查脚本克隆后对 `.git` 的处理，确认当前实现与需求不符。
- 根因：README 未同步模板仓库实际地址；脚本设计成“完全脱离模板”，直接删除 git 元数据。
- 修复建议：修正文档地址，保留 `.git`，将 `origin` 重命名为 `template` 并统一为 HTTPS。
- 回归验证点：初始化后存在 `.git`，`git remote -v` 仅显示 `template`，生成项目 README 不再残留模板说明。

### 2. admin 首次安装存在 npm 审计告警

- 症状：初始化过程中执行 `npm install` 后会出现 1 个 moderate 和 1 个 high 漏洞提示。
- 复现步骤：在 `admin` 目录运行 `npm audit --json`。
- 影响范围：所有使用模板初始化并安装前端依赖的项目。
- 假设与排查过程：通过 `npm audit --json`、`npm ls brace-expansion picomatch` 和 `npm audit fix --dry-run --json` 确认告警来自 `brace-expansion@2.0.2` 与 `picomatch@2.3.1/4.0.3`。
- 根因：模板仓库锁文件中仍解析到已知漏洞版本的间接依赖。
- 修复建议：刷新 `admin/package-lock.json` 到 `brace-expansion@2.0.3`、`picomatch@2.3.2/4.0.4`。
- 回归验证点：生成项目首次安装完成后执行 `npm audit --json`，漏洞数为 0。
