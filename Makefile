APP_NAME := gin-template
GO_TOOL_ENV := GOPROXY=https://proxy.golang.org,direct GOCACHE=$(CURDIR)/.cache/go/build GOMODCACHE=$(CURDIR)/.cache/go/mod GOLANGCI_LINT_CACHE=$(CURDIR)/.cache/go/golangci-lint
LOCAL_BIN := $(CURDIR)/bin

.PHONY: setup dev dev-backend dev-admin build lint test admin-check check generate verify ci-local tidy contract-fuzz

dev:
	@echo "请分别在两个终端执行 'make dev-backend' 和 'make dev-admin'"

setup:
	env $(GO_TOOL_ENV) go mod tidy
	env $(GO_TOOL_ENV) GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2
	env $(GO_TOOL_ENV) GOBIN=$(LOCAL_BIN) go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.0
	cd admin && npm install
	@if command -v lefthook >/dev/null 2>&1; then lefthook install; else echo "skip lefthook install (command not found)"; fi

dev-backend:
	env $(GO_TOOL_ENV) go run ./cmd/server serve --env dev --config configs/config.dev.yaml

dev-admin:
	cd admin && npm run dev

build:
	cd admin && npm run build
	env $(GO_TOOL_ENV) go build -o bin/$(APP_NAME) ./cmd/server

lint:
	env $(GO_TOOL_ENV) go vet ./...
	@env $(GO_TOOL_ENV) sh -c 'if [ -x "$(LOCAL_BIN)/golangci-lint" ]; then "$(LOCAL_BIN)/golangci-lint" run; elif command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; else echo "golangci-lint 未安装，请先执行 make setup" && exit 1; fi'
	cd admin && npm run lint

test:
	env $(GO_TOOL_ENV) go test ./...
	cd admin && npm run test

admin-check:
	cd admin && npm run lint
	cd admin && npm run typecheck
	cd admin && npm run test
	cd admin && npm run build

check:
	env $(GO_TOOL_ENV) go run ./cmd/server config verify --env dev
	env $(GO_TOOL_ENV) go run ./cmd/server openapi validate
	make lint
	make test
	cd admin && npm run typecheck
	cd admin && npm run build

gen:
	env $(GO_TOOL_ENV) go run ./cmd/server openapi generate

verify: gen
	@test -f internal/app/http/openapi.gen.go
	@test -f admin/src/api/generated.ts
	@git diff --exit-code -- internal/app/http/openapi.gen.go admin/src/api/generated.ts

ci-local:
	make verify
	make check

tidy:
	env $(GO_TOOL_ENV) go mod tidy

contract-fuzz:
	@if ! command -v st >/dev/null 2>&1 && ! command -v schemathesis >/dev/null 2>&1; then echo "schemathesis 未安装，请先安装 st 或 schemathesis 命令" && exit 1; fi
	@tmpdir=$$(mktemp -d); \
	schemathesis_cmd=$$(if command -v st >/dev/null 2>&1; then printf %s st; else printf %s schemathesis; fi); \
	db_path="$$tmpdir/contract-fuzz.db"; \
	upload_dir="$$tmpdir/uploads"; \
	base_url="http://127.0.0.1:3301"; \
	env $(GO_TOOL_ENV) APP_STORAGE_LOCAL_DIR="$$upload_dir" APP_RATE_LIMIT_ENABLED=false APP_TURNSTILE_ENABLED=false APP_MAIL_ENABLED=false go run ./cmd/server serve --env dev --config configs/config.dev.yaml --host 127.0.0.1 --port 3301 --sqlite-path "$$db_path" >"$$tmpdir/server.log" 2>&1 & \
	pid=$$!; \
	trap 'kill $$pid >/dev/null 2>&1 || true; wait $$pid >/dev/null 2>&1 || true; rm -rf "$$tmpdir"' EXIT; \
	attempt=0; \
	until curl -sf "$$base_url/api/v1/system/status" >/dev/null; do \
		attempt=$$((attempt + 1)); \
		if [ $$attempt -ge 30 ]; then \
			cat "$$tmpdir/server.log"; \
			exit 1; \
		fi; \
		sleep 1; \
	done; \
	"$$schemathesis_cmd" run api/openapi/openapi.yaml --url "$$base_url" --phases examples,coverage,fuzzing --checks not_a_server_error,status_code_conformance,content_type_conformance,response_schema_conformance
