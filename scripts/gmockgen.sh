#!/bin/bash
# 概要
#   usecase / repository インターフェースの gomock を生成します。
#   各 .go ファイルに定義された //go:generate ディレクティブを go generate で実行します。
#
# 使用方法（プロジェクトルートから実行）
#   ./scripts/gmockgen.sh
#
# オプション
#   -u, --usecase      usecase のみ生成
#   -r, --repository   repository のみ生成
#   -h, --help         ヘルプを表示
#
# 例
#   ./scripts/gmockgen.sh
#   ./scripts/gmockgen.sh --usecase
#   ./scripts/gmockgen.sh --repository
#
# 対象パッケージ
#   internal/adapter/interface/usecase/...
#   internal/domain/interface/repository/...
#
# 注意
#   - プロジェクトルート（go.mod があるディレクトリ）から実行してください。
#   - mockgen は go tool mockgen 経由で実行されます（go.mod の tool 依存を使用）。

set -euo pipefail

GENERATE_USECASE=true
GENERATE_REPOSITORY=true

usage() {
    sed -n '2,20p' "$0" | sed 's/^# \{0,1\}//'
}

while [[ "$#" -gt 0 ]]; do
    case "$1" in
        -u|--usecase)
            GENERATE_USECASE=true
            GENERATE_REPOSITORY=false
            shift
            ;;
        -r|--repository)
            GENERATE_USECASE=false
            GENERATE_REPOSITORY=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown parameter passed: $1" >&2
            usage
            exit 1
            ;;
    esac
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

if [[ ! -f "${PROJECT_ROOT}/go.mod" ]]; then
    echo "Error: go.mod not found. Run this script from the project root." >&2
    exit 1
fi

cd "${PROJECT_ROOT}"

echo "=== Generating gomock mocks ==="
echo "Project root: ${PROJECT_ROOT}"

if [[ "${GENERATE_USECASE}" == true ]]; then
    echo
    echo "[usecase] go generate ./internal/adapter/interface/usecase/..."
    go generate ./internal/adapter/interface/usecase/...
fi

if [[ "${GENERATE_REPOSITORY}" == true ]]; then
    echo
    echo "[repository] go generate ./internal/domain/interface/repository/..."
    go generate ./internal/domain/interface/repository/...
fi

echo
echo "=== Done ==="
