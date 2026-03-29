# 50 Test Report

## 触达范围

- README 初始化文档
- `scripts/init-project.sh` 的 git 远程、README 清理与安装流程
- `admin/package-lock.json` 的间接依赖安全版本
- 仓库标准门禁命令

## 执行命令

- `bash -n scripts/init-project.sh`
- `bash scripts/init-project.sh --module github.com/acme/tearsapp --target /tmp/gin-template-init-noinstall --repo-url git@github.com:fengjx/gin-template.git --skip-install`
- `git -C /tmp/gin-template-init-noinstall remote -v`
- `test -d /tmp/gin-template-init-noinstall/.git`
- `rm -rf /tmp/gin-template-template-worktree && rsync -a --delete --exclude '.git' --exclude 'admin/node_modules' --exclude '.cache' --exclude 'runtime' ./ /tmp/gin-template-template-worktree/ && git -C /tmp/gin-template-template-worktree init && git -C /tmp/gin-template-template-worktree add . && git -C /tmp/gin-template-template-worktree -c user.name=Codex -c user.email=codex@example.com commit -m 'temp template snapshot'`
- `bash scripts/init-project.sh --module github.com/acme/tearsapp --target /tmp/gin-template-init-local --repo-url /tmp/gin-template-template-worktree --ref master --skip-install`
- `rg -n "template:init|raw.githubusercontent|模板初始化脚本|git fetch template|模板仓库地址" /tmp/gin-template-init-local/README.md`
- `bash scripts/init-project.sh --module github.com/acme/tearsapp --target /tmp/gin-template-init-full --repo-url /tmp/gin-template-template-worktree --ref master`
- `cd /tmp/gin-template-init-full/admin && npm audit --json`
- `make gen`
- `make verify`
- `make lint`
- `make test`
- `make check`

## 结果

- `bash -n scripts/init-project.sh` 通过。
- GitHub SSH 地址回归中，生成项目保留 `.git`，且 `git remote -v` 仅显示 `template https://github.com/fengjx/gin-template.git`。
- 基于当前工作树快照模板仓库的无安装回归中，生成项目 README 已不再包含模板初始化说明块。
- 基于当前工作树快照模板仓库的完整初始化中，`npm install` 完成后 `npm audit --json` 返回 `total=0`。
- `make gen`、`make verify`、`make lint`、`make test`、`make check` 全部通过。

## 未覆盖项

- 未针对 GitHub Enterprise 或其他自建 Git 服务的 SSH 地址做归一化验证；当前实现按设计仅处理 `github.com`。
- 未验证用户后续手动新增 `origin` 远程的流程；本次需求明确不自动创建该远程。

## 结论

- `通过`
