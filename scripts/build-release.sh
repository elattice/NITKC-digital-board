#!/bin/sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
PROJECT_ROOT=$(dirname "$SCRIPT_DIR")
DIST_DIR="$PROJECT_ROOT/dist"

command -v npm >/dev/null 2>&1 || {
  echo "error: npm is required" >&2
  exit 1
}
command -v go >/dev/null 2>&1 || {
  echo "error: go is required" >&2
  exit 1
}

echo "Building embedded React frontend..."
(cd "$PROJECT_ROOT/frontend" && npm run build)

mkdir -p "$DIST_DIR"

echo "Building Apple Silicon macOS binary..."
(cd "$PROJECT_ROOT/backend" && env \
  CGO_ENABLED=0 \
  GOOS=darwin \
  GOARCH=arm64 \
  go build -trimpath -o "$DIST_DIR/digital-board-macos-arm64" ./cmd/server)

echo "Building Intel/AMD Linux binary..."
(cd "$PROJECT_ROOT/backend" && env \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  go build -trimpath -o "$DIST_DIR/digital-board-linux-amd64" ./cmd/server)

echo "Build complete:"
echo "  $DIST_DIR/digital-board-macos-arm64"
echo "  $DIST_DIR/digital-board-linux-amd64"
