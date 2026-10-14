#!/bin/sh
# GitForge installation script with SHA-256 checksum verification
set -eu

REPO="phantom1785322-hub/gitforge"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# 1. Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    linux*)   TARGET_OS="linux" ;;
    darwin*)  TARGET_OS="darwin" ;;
    freebsd*) TARGET_OS="freebsd" ;;
    openbsd*) TARGET_OS="openbsd" ;;
    *)
        echo "Error: Unsupported operating system $OS" >&2
        exit 1
        ;;
esac

# 2. Detect Architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64) TARGET_ARCH="amd64" ;;
    arm64|aarch64) TARGET_ARCH="arm64" ;;
    *)
        echo "Error: Unsupported architecture $ARCH" >&2
        exit 1
        ;;
esac

echo "Detected platform: ${TARGET_OS}/${TARGET_ARCH}"

VERSION="${1:-latest}"
if [ "$VERSION" = "latest" ]; then
    VERSION_TAG="v0.1.0"
else
    VERSION_TAG="$VERSION"
fi

TARBALL="gitforge_${VERSION_TAG#v}_${TARGET_OS}_${TARGET_ARCH}.tar.gz"
BASE_URL="https://github.com/${REPO}/releases/download/${VERSION_TAG}"

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading ${TARBALL}..."
curl -sSL "${BASE_URL}/${TARBALL}" -o "${TMPDIR}/${TARBALL}"
curl -sSL "${BASE_URL}/checksums.txt" -o "${TMPDIR}/checksums.txt"

echo "Verifying SHA-256 checksum..."
cd "$TMPDIR"
if command -v sha256sum >/dev/null 2>&1; then
    grep "$TARBALL" checksums.txt | sha256sum -c -
elif command -v shasum >/dev/null 2>&1; then
    grep "$TARBALL" checksums.txt | shasum -a 256 -c -
else
    echo "Warning: neither sha256sum nor shasum found; cannot verify checksum." >&2
fi

echo "Extracting binary..."
tar -xzf "$TARBALL"

if [ ! -w "$INSTALL_DIR" ]; then
    echo "Installing to $INSTALL_DIR (requires sudo)..."
    sudo install -m 755 gitforge "$INSTALL_DIR/gitforge"
else
    echo "Installing to $INSTALL_DIR..."
    install -m 755 gitforge "$INSTALL_DIR/gitforge"
fi

echo "GitForge installed successfully:"
"$INSTALL_DIR/gitforge" version
