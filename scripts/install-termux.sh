#!/data/data/com.termux/files/usr/bin/bash
# GitForge native build and install script for Termux (Android)
set -euo pipefail

echo "=== GitForge Termux Installer ==="

# 1. Ensure required dependencies
pkg update -y
pkg install -y git golang make

# 2. Verify Go environment
echo "Go version: $(go version)"

# 3. Build GitForge natively with CGO_ENABLED=0
export CGO_ENABLED=0
echo "Building GitForge static binary..."
go build -ldflags "-s -w -X github.com/phantom1785322-hub/gitforge/internal/core.BuildVersion=termux" \
  -o "$PREFIX/bin/gitforge" ./cmd/gitforge

# 4. Verification
echo "Verifying installation..."
gitforge version
gitforge doctor

echo "=== GitForge installed successfully to $PREFIX/bin/gitforge ==="
