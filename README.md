# gin-template

一个基于 Gin + React + Vite 的前后端同构项目脚手架，默认支持 SQLite，按需切换 MySQL，OpenAPI 契约驱动前后端接口与文档。

## 特性

- Gin 后端 + React 19 前端 + Vite 8 构建
- Viper 配置加载，支持配置文件、环境变量和 CLI 参数
- SQLite 默认回退，MySQL 可选
- 新库自动执行 bootstrap SQL，后续 schema 通过手工 SQL 升级
- JWT Access Token + Refresh Token Cookie
- Zap `applog` / `accesslog` 分离
- `trace_id` 全链路透传
- OpenAPI 文档与前端类型生成脚本
- 面向 AI 编码助手的仓库宪法、闭环流程与 feature 工作区
- `make check` / `make verify` 一体化自检链路
- 多层配置加载与 `.env.example` 环境变量示例
- React 管理台：登录、用户、文件、配置、个人信息、文档页

日志输出约定：

- 运行环境由 `APP_ENV` 或 `--env` 指定，只支持 `dev`、`test`、`prod`，默认 `dev`
- `dev` 环境下：`applog` 输出控制台，`accesslog` 写入 `runtime/logs/access.log`
- 非 `dev` 环境下：`applog` 写入 `runtime/logs/app.log`，`accesslog` 写入 `runtime/logs/access.log`
- 文件日志按自然天切割，文件名会自动追加日期后缀，例如 `app-2026-03-22.log`

## 快速开始

```bash
make setup
make dev-backend
# 新开一个终端
make dev-admin
```

默认地址：

- 后端: `http://localhost:3000`
- 前端: `http://localhost:5173`
- 文档: `http://localhost:3000/docs`
- 首次启动空库时会自动创建管理员账号: `admin / admin`

## 常用命令

```bash
go run ./cmd/server config verify --env dev
go run ./cmd/server serve --config configs/config.dev.yaml
go run ./cmd/server schema verify --config configs/config.dev.yaml
go run ./cmd/server openapi validate
go run ./cmd/server seed admin --password your-password
make check
make verify
make build
make test
```

## 统一开发入口

```bash
make setup
make gen
make verify
make lint
make test
make check
make ci-local
```

- `make setup`：安装 Go / Admin 依赖，并在本机有 `lefthook` 时安装 git hooks
- `make gen`：生成 Go 服务端 OpenAPI 文件和前端 TS 类型
- `make verify`：生成后检查生成物是否漂移
- `make check`：配置校验、OpenAPI 校验、lint、test、前端 typecheck/build
- `make ci-local`：本地执行与 CI 接近的一套完整检查

## 配置加载规则

业务配置优先级从低到高如下：

1. 代码默认值
2. `configs/config.yaml`
3. `configs/config.{env}.yaml`
4. `configs/config.local.yaml`
5. `APP_` 环境变量
6. CLI 参数

运行环境优先级从低到高如下：

1. 默认值 `dev`
2. `APP_ENV`
3. `--env`

示例：

```bash
APP_SERVER_PORT=4000 APP_AUTH_JWT_SECRET=replace-me go run ./cmd/server serve --env dev
```

环境变量与配置项的映射规则：

- `server.port` -> `APP_SERVER_PORT`
- `auth.jwt_secret` -> `APP_AUTH_JWT_SECRET`
- `database.sqlite_path` -> `APP_DATABASE_SQLITE_PATH`

敏感信息请放在本地环境变量或 `configs/config.local.yaml`，参考根目录 `.env.example`。

## 数据库规则

- 当 `database.dsn` 为空时，自动使用 SQLite
- SQLite 默认文件: `runtime/data/app.db`
- 新库首次启动时自动执行 `database/bootstrap/sqlite.sql` 或 `database/bootstrap/mysql.sql`
- 后续改库请手工执行 `database/manual/sqlite/*.sql` 或 `database/manual/mysql/*.sql`

## 前端类型生成

```bash
go run ./cmd/server openapi generate
```

这条命令会调用 `oapi-codegen` 和 `openapi-typescript`，分别生成 Go 服务端类型和前端 TS 类型。

## AI 协作文件

- `AGENTS.md`：仓库规范唯一来源
- `docs/architecture/README.md`：架构速览
- `docs/playbooks/feature-lifecycle.md`：AI 闭环流程入口
- `docs/features/README.md`：feature 工作区说明与模板约定
- `docs/code_review.md`：统一 review 基线
- `docs/playbooks/`：新增接口、加表、加页面、加配置的原子剧本
- `.agents/skills/`：仓库共享 skills 唯一权威目录

## AI 闭环工作流

标准 feature：

`flow-idea -> flow-prd -> flow-spec -> flow-tasks -> flow-impl -> flow-test + flow-review -> flow-investigate(如有问题) -> flow-doc-release`

bugfix 快路径：

`flow-investigate -> flow-spec -> flow-tasks -> flow-impl -> flow-test -> flow-review -> flow-doc-release`

每个非微小功能改动都需要创建 `docs/features/<NNN>-<slug>/`，并把过程文档按阶段沉淀到目录中。

`docs/features/<NNN>-<slug>/` 同时也是本仓库的 execution plan 载体，不再额外维护独立 `PLANS.md`。

<!-- template:init:start -->
## 模板初始化脚本

如果这个仓库已经发布到 GitHub，可以直接用 `curl` 一键拉到本地并替换成自己的 module：

```bash
curl -fsSL https://raw.githubusercontent.com/fengjianxin/gin-template/main/scripts/init-project.sh | bash
```

无交互模式：

```bash
curl -fsSL https://raw.githubusercontent.com/fengjianxin/gin-template/main/scripts/init-project.sh | \
  bash -s -- --module github.com/acme/my-app --target my-app
```

脚本会完成这些事：

- 克隆模板到目标目录
- 把默认 module `gin-template` 替换成你输入的自定义 module
- 如果你传入的是完整 module，例如 `github.com/acme/my-app`，会自动提取短名 `my-app` 作为应用名
- 同步替换 Go import、应用名、前端包名、Docker 二进制名
- 可选执行 `go mod tidy` 和 `npm install`
