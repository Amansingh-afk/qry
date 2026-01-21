#!/bin/sh
set -e

INSTALL_DIR="/usr/local/bin"

# Colors
CYAN='\033[36m'
PINK='\033[35m'
GREEN='\033[32m'
YELLOW='\033[33m'
RESET='\033[0m'

printf "\n"
printf "  ${CYAN}Q${PINK}R${CYAN}Y${RESET} Uninstaller\n"
printf "\n"

# Check if qry exists
if [ ! -f "$INSTALL_DIR/qry" ]; then
    # Try to find it
    QRY_PATH=$(which qry 2>/dev/null || true)
    if [ -z "$QRY_PATH" ]; then
        printf "  ${YELLOW}!${RESET} qry is not installed\n"
        printf "\n"
        exit 0
    fi
    INSTALL_DIR=$(dirname "$QRY_PATH")
fi

printf "  ${CYAN}→${RESET} Removing %s/qry\n" "$INSTALL_DIR"

if [ -w "$INSTALL_DIR/qry" ]; then
    rm "$INSTALL_DIR/qry"
else
    sudo rm "$INSTALL_DIR/qry"
fi

printf "\n"
printf "  ${GREEN}✓${RESET} qry uninstalled successfully\n"
printf "\n"
printf "  ${YELLOW}Note:${RESET} Project configs (.qry.yaml, .qry/) were not removed.\n"
printf "        Remove them manually if needed.\n"
printf "\n"
