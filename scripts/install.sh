#!/bin/sh
set -e

REPO="amansingh-afk/qry"
INSTALL_DIR="/usr/local/bin"

echo ""
echo "  \033[36mQ\033[35mR\033[36mY\033[0m Installer"
echo ""

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    linux|darwin) ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

LATEST=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)

if [ -z "$LATEST" ]; then
    echo "Could not fetch latest release"
    exit 1
fi

URL="https://github.com/$REPO/releases/download/$LATEST/qry_${OS}_${ARCH}.tar.gz"

echo "  → Downloading qry $LATEST for $OS/$ARCH"

TMP=$(mktemp -d)
curl -fsSL "$URL" | tar -xz -C "$TMP"

echo "  → Installing to $INSTALL_DIR"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP/qry" "$INSTALL_DIR/qry"
else
    sudo mv "$TMP/qry" "$INSTALL_DIR/qry"
fi

rm -rf "$TMP"

echo ""
echo "  \033[32m✓\033[0m qry installed successfully"
echo ""
echo "  Run: qry --help"
echo ""
