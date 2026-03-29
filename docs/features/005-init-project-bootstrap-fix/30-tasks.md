# 30 Tasks

## 任务列表

### T1

- 目标：修复 README 的初始化入口和模板说明块约定。
- 修改范围：`README.md`
- 前置依赖：`20-tech-spec.md`
- 产出物：正确的 raw 地址、模板仓库地址、`template:init` 起止标记和初始化后 git 说明
- 验证命令：`rg -n "raw.githubusercontent.com/fengjx/gin-template|template:init|模板仓库|git fetch template" README.md`

### T2

- 目标：调整初始化脚本保留历史并整理模板远程。
- 修改范围：`scripts/init-project.sh`
- 前置依赖：T1
- 产出物：保留 `.git`、将 `origin` 重命名为 `template`、GitHub 地址 HTTPS 归一化、完成提示更新
- 验证命令：`bash scripts/init-project.sh --module github.com/acme/tearsapp --target /tmp/gin-template-init-noinstall --skip-install`

### T3

- 目标：消除 admin 初始化依赖的已知审计告警。
- 修改范围：`admin/package-lock.json`
- 前置依赖：T2
- 产出物：升级到安全版本的 lockfile
- 验证命令：`cd admin && npm audit --json`

### T4

- 目标：完成 bugfix 闭环文档、测试、review、提交流程和 PR 记录。
- 修改范围：`docs/features/005-init-project-bootstrap-fix/*`
- 前置依赖：T1-T3
- 产出物：实现记录、测试报告、review 结论、发布说明、PR 链接
- 验证命令：`make gen && make verify && make lint && make test && make check`
