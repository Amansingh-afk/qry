#!/bin/sh
set -e

REPO="amansingh-afk/qry"
INSTALL_DIR="/usr/local/bin"

# Colors
CYAN='\033[36m'
PINK='\033[35m'
GREEN='\033[32m'
RESET='\033[0m'

printf "\n"
printf "  ${CYAN}Q${PINK}R${CYAN}Y${RESET} Installer\n"
printf "\n"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) printf "  Unsupported architecture: %s\n" "$ARCH"; exit 1 ;;
esac

case "$OS" in
    linux|darwin) ;;
    *) printf "  Unsupported OS: %s\n" "$OS"; exit 1 ;;
esac

LATEST=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)

if [ -z "$LATEST" ]; then
    printf "  Could not fetch latest release\n"
    exit 1
fi

URL="https://github.com/$REPO/releases/download/$LATEST/qry_${OS}_${ARCH}.tar.gz"

printf "  ${CYAN}→${RESET} Downloading qry %s for %s/%s\n" "$LATEST" "$OS" "$ARCH"

TMP=$(mktemp -d)
curl -fsSL "$URL" | tar -xz -C "$TMP"

printf "  ${CYAN}→${RESET} Installing to %s\n" "$INSTALL_DIR"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP/qry" "$INSTALL_DIR/qry"
else
    sudo mv "$TMP/qry" "$INSTALL_DIR/qry"
fi

rm -rf "$TMP"

printf "\n"
printf "  ${GREEN}✓${RESET} qry installed successfully\n"
printf "\n"
printf "  Run: ${CYAN}qry --help${RESET}\n"
printf "\n"
