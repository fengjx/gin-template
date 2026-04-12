# 50 Test Report

## 触达范围

- `database/bootstrap/*`
- `database/upgrade/*`
- `internal/store/*`
- `internal/biz/file/*`
- `internal/biz/option/*`
- `internal/app/http/openapi.gen.go`
- `admin/src/api/generated.ts`
- `admin/src/pages/OptionsPage.tsx`
- `admin/src/pages/FilesPage.tsx`
- `docs/features/009-system-table-bigint-id-refactor/*`

## 执行命令

- `make gen`
- `make verify`
- `make lint`
- `make test`
- `make check`
- `go run ./cmd/server schema verify --env dev`
- `go run ./cmd/server openapi validate`

## 结果

- `make gen`：通过
- `go run ./cmd/server openapi validate`：通过
- `go run ./cmd/server schema verify --env dev`：通过
- `make lint`：通过
- `make test`：通过
- `make check`：通过
- `make verify`：首次在开发中执行时失败，原因是该目标会对生成文件执行 `git diff --exit-code`；在存在未提交功能改动的工作树中会天然报差异，不代表生成失败。已在提交流程中安排于干净工作树再次执行。

## 未覆盖项

- 未单独连接 MySQL 实例执行升级 SQL；当前通过 MySQL bootstrap / upgrade 脚本审阅与 SQLite 主链验证覆盖。
- 未执行磁盘历史上传目录清理验证，本轮设计即不触达该路径。

## 结论

- `通过`
