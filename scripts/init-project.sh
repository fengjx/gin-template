#!/usr/bin/env bash
set -euo pipefail

DEFAULT_MODULE="gin-template"
DEFAULT_REPO_URL="${TEMPLATE_REPO_URL_DEFAULT:-https://github.com/fengjx/gin-template.git}"
DEFAULT_REF="${TEMPLATE_REPO_REF_DEFAULT:-main}"

MODULE_NAME=""
TARGET_DIR=""
REPO_URL="$DEFAULT_REPO_URL"
REPO_REF="$DEFAULT_REF"
SKIP_INSTALL="false"

usage() {
  cat <<'EOF'
用法:
  init-project.sh [--module <go-module>] [--target <dir>] [--repo-url <git-url>] [--ref <git-ref>] [--skip-install]

示例:
  curl -fsSL <raw-script-url> | bash
  curl -fsSL <raw-script-url> | bash -s -- --module github.com/acme/my-app --target my-app
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --module)
      MODULE_NAME="${2:-}"
      shift 2
      ;;
    --target)
      TARGET_DIR="${2:-}"
      shift 2
      ;;
    --repo-url)
      REPO_URL="${2:-}"
      shift 2
      ;;
    --ref)
      REPO_REF="${2:-}"
      shift 2
      ;;
    --skip-install)
      SKIP_INSTALL="true"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "未知参数: $1" >&2
      usage
      exit 1
      ;;
  esac
done

prompt_if_empty() {
  local prompt="$1"
  local default_value="$2"
  local result
  read -r -p "${prompt} [${default_value}]: " result
  if [[ -z "$result" ]]; then
    echo "$default_value"
  else
    echo "$result"
  fi
}

if [[ -z "$MODULE_NAME" ]]; then
  MODULE_NAME="$(prompt_if_empty "请输入新的 Go module 名称" "$DEFAULT_MODULE")"
fi

APP_NAME="${MODULE_NAME##*/}"
if [[ -z "$APP_NAME" ]]; then
  APP_NAME="$MODULE_NAME"
fi

if [[ -z "$TARGET_DIR" ]]; then
  TARGET_DIR="$(prompt_if_empty "请输入目标目录" "$APP_NAME")"
fi

if [[ -e "$TARGET_DIR" ]] && [[ -n "$(find "$TARGET_DIR" -mindepth 1 -maxdepth 1 2>/dev/null)" ]]; then
  echo "目标目录 '$TARGET_DIR' 已存在且非空，请换一个目录。" >&2
  exit 1
fi

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "缺少依赖命令: $1" >&2
    exit 1
  fi
}

need_cmd git
need_cmd perl

echo "==> 克隆模板仓库"
git clone --depth=1 --branch "$REPO_REF" "$REPO_URL" "$TARGET_DIR"
rm -rf "$TARGET_DIR/.git"
rm -rf "$TARGET_DIR/.idea" "$TARGET_DIR/.cache" "$TARGET_DIR/runtime" "$TARGET_DIR/admin/node_modules" "$TARGET_DIR/admin/dist/assets" "$TARGET_DIR/admin/tsconfig.tsbuildinfo"

replace_literal() {
  local needle="$1"
  local replacement="$2"
  shift 2
  if [[ $# -eq 0 ]]; then
    return 0
  fi
  local existing_files=()
  local file
  for file in "$@"; do
    if [[ -f "$file" ]]; then
      existing_files+=("$file")
    fi
  done
  if [[ ${#existing_files[@]} -eq 0 ]]; then
    return 0
  fi
  perl -0pi -e "s#\Q${needle}\E#${replacement}#g" "${existing_files[@]}"
}

echo "==> 替换 Go module 和内部 import"
GO_FILES=()
while IFS= read -r -d '' file; do
  GO_FILES+=("$file")
done < <(find "$TARGET_DIR" -type f \( -name '*.go' -o -name 'go.mod' \) -print0)
if [[ ${#GO_FILES[@]} -gt 0 ]]; then
  replace_literal "$DEFAULT_MODULE" "$MODULE_NAME" "${GO_FILES[@]}"
fi
replace_literal "Use:   \"$MODULE_NAME\"" "Use:   \"$APP_NAME\"" "$TARGET_DIR/internal/app/command/root.go"
replace_literal "v.SetDefault(\"app.name\", \"$MODULE_NAME\")" "v.SetDefault(\"app.name\", \"$APP_NAME\")" "$TARGET_DIR/internal/app/config/config.go"
replace_literal "v.SetDefault(\"auth.issuer\", \"$MODULE_NAME\")" "v.SetDefault(\"auth.issuer\", \"$APP_NAME\")" "$TARGET_DIR/internal/app/config/config.go"

echo "==> 替换应用名称"
TEXT_FILES=()
while IFS= read -r -d '' file; do
  TEXT_FILES+=("$file")
done < <(find "$TARGET_DIR" -type f \
  \( -name 'README.md' -o -name 'Makefile' -o -name 'Dockerfile' -o -name '*.yaml' -o -name '*.yml' -o -name '*.toml' -o -name '*.md' -o -name '*.mdc' -o -name '.env.example' -o -name '*.json' -o -name 'index.html' -o -name '*.ts' -o -name '*.tsx' \) \
  -not -path '*/node_modules/*' -not -path '*/.git/*' -print0)

replace_literal "APP_NAME := gin-template" "APP_NAME := ${APP_NAME}" "$TARGET_DIR/Makefile"
replace_literal "\"name\": \"gin-template-admin\"" "\"name\": \"${APP_NAME}-admin\"" "$TARGET_DIR/admin/package.json" "$TARGET_DIR/admin/package-lock.json"
replace_literal "title: gin-template API" "title: ${APP_NAME} API" "$TARGET_DIR/api/openapi/openapi.yaml"
replace_literal "> gin-template<" "> ${APP_NAME} <" "$TARGET_DIR/admin/index.html" "$TARGET_DIR/admin/dist/index.html"
if [[ ${#TEXT_FILES[@]} -gt 0 ]]; then
  replace_literal "gin-template" "$APP_NAME" "${TEXT_FILES[@]}"
fi
perl -0pi -e 's@\n?<!-- template:init:start -->.*?<!-- template:init:end -->\n?@\n@s' "$TARGET_DIR/README.md" 2>/dev/null || true

echo "==> 清理构建产物"
find "$TARGET_DIR" -name '.DS_Store' -delete || true

if [[ "$SKIP_INSTALL" == "false" ]]; then
  if command -v go >/dev/null 2>&1; then
    echo "==> 执行 go mod tidy"
    (cd "$TARGET_DIR" && go mod tidy)
  else
    echo "==> 未检测到 go，跳过 go mod tidy"
  fi

  if command -v npm >/dev/null 2>&1; then
    echo "==> 执行 npm install"
    (cd "$TARGET_DIR/admin" && npm install)
  else
    echo "==> 未检测到 npm，跳过 npm install"
  fi
fi

cat <<EOF

初始化完成
  目录:   $TARGET_DIR
  module: $MODULE_NAME
  app:    $APP_NAME

建议下一步:
  cd $TARGET_DIR
  make dev-backend
  make dev-admin
EOF
