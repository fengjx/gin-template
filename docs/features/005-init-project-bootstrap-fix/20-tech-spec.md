# 20 Tech Spec

## 目标与非目标

- 目标：修复 README 初始化入口、让初始化项目保留模板历史并显式保留 `template` 远程、消除前端初始化依赖告警。
- 非目标：新增初始化参数、修改 API/配置/数据库契约、替用户自动绑定其私有仓库远程。

## 现状与约束

- README 初始化命令仍指向错误仓库路径，且模板区块只有起始注释，没有结束注释。
- `scripts/init-project.sh` 当前克隆后直接删除 `.git`，与“保留模板同步能力”冲突。
- 仓库约束要求默认使用 HTTPS 模板地址，并把重要改动沉淀到 `docs/features/<NNN>-<slug>/`。
- npm 审计显示问题集中在 `brace-expansion` 和 `picomatch` 的间接依赖版本偏低，可通过刷新 lockfile 修复。

## 模块边界

- `README.md`：修复公开初始化文档，补充模板仓库地址与初始化后 git 同步说明。
- `scripts/init-project.sh`：负责模板克隆、git 远程整理、文本替换与安装步骤。
- `admin/package-lock.json`：记录前端依赖解析结果，确保新项目安装落到安全版本。
- `docs/features/005-init-project-bootstrap-fix/*`：记录调查、实现、测试、review 与发布闭环。

## 接口与数据流

- 初始化命令入口不变：`curl .../scripts/init-project.sh | bash`。
- 脚本流程调整为：
  - 克隆模板仓库到目标目录。
  - 归一化 GitHub 模板地址为 HTTPS。
  - 保留 `.git`，将 `origin` 重命名为 `template` 并写回归一化后的地址。
  - 执行项目名、module、README 等文本替换。
  - 可选执行 `go mod tidy` 与 `npm install`。
- `--repo-url` 输入规则：
  - `git@github.com:owner/repo.git` -> `https://github.com/owner/repo.git`
  - `ssh://git@github.com/owner/repo.git` -> `https://github.com/owner/repo.git`
  - 已经是 `https://github.com/...` 的地址保持不变
  - 非 GitHub 地址保持原样

## OpenAPI / 配置 / 数据库影响

- 无 OpenAPI 变更。
- 无新增配置项。
- 无数据库 schema 变更。

## 风险、回滚与观测点

- 风险：脚本误处理非 GitHub 地址或重复重命名远程，导致初始化仓库远程异常。
- 风险：lockfile 升级可能引入前端构建回归。
- 回滚：恢复 README、脚本和 `admin/package-lock.json` 到当前主线版本即可。
- 观测点：初始化后 `git remote -v`、生成项目 README 内容、`npm audit --json` 结果、前端与全仓门禁命令。

## 测试与验证策略

- 静态检查 README 中 raw 地址、模板仓库地址与模板起止标记。
- 在临时目录运行 `bash scripts/init-project.sh --skip-install`，验证保留 `.git`、远程为 `template`、地址为 HTTPS、README 模板区块已删除。
- 在临时目录运行完整初始化并于生成项目 `admin` 执行 `npm audit --json`，确认漏洞数为 0。
- 执行 `make gen`、`make verify`、`make lint`、`make test`、`make check`。
