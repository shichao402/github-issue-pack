#!/bin/bash
# 构建脚本 - macOS/Linux 入口

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

echo "=== GitHub Issue Pack 构建 ==="

# 读取版本
VERSION=$(grep -o '"version": "[^"]*"' version.json | head -1 | cut -d'"' -f4)
echo "版本: $VERSION"

# 创建输出目录
mkdir -p bin

# 构建目标平台
PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS="${PLATFORM%/*}"
    GOARCH="${PLATFORM#*/}"
    OUTPUT="bin/github-issue-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi
    
    echo "构建 $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w -X github.com/shichao402/github-issue-pack/internal/cli.Version=$VERSION" -o "$OUTPUT" ./cmd/github-issue
done

# 创建当前平台的符号链接
CURRENT_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
CURRENT_ARCH=$(uname -m)
if [ "$CURRENT_ARCH" = "x86_64" ]; then
    CURRENT_ARCH="amd64"
elif [ "$CURRENT_ARCH" = "aarch64" ] || [ "$CURRENT_ARCH" = "arm64" ]; then
    CURRENT_ARCH="arm64"
fi

ln -sf "github-issue-${CURRENT_OS}-${CURRENT_ARCH}" bin/github-issue

echo ""
echo "✅ 构建完成!"
ls -la bin/
