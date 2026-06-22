#!/bin/bash
set -e

# ═══════════════════════════════════════════════
#   Term1zz — Uninstaller
#   https://codeberg.org/Den1zz/Term1zz
# ═══════════════════════════════════════════════

INSTALL_DIR="$HOME/.local/share/Term1zz"

# ── Colors ──────────────────────────────────────
BOLD='\033[1m'
DIM='\033[2m'
GREEN='\033[38;2;166;227;161m'
RED='\033[38;2;243;139;168m'
BLUE='\033[38;2;137;180;250m'
RESET='\033[0m'

info()  { echo -e "${BLUE}${BOLD}  ▸${RESET} $1"; }
ok()    { echo -e "${GREEN}${BOLD}  ✓${RESET} $1"; }
fail()  { echo -e "${RED}${BOLD}  ✗${RESET} $1"; exit 1; }

echo ""
echo -e "${RED}${BOLD}  ⚠ WARNING: Uninstallation${RESET}"
echo -e "  This will remove all Term1zz symlinks and attempt to restore"
echo -e "  your latest configurations from the backup directory."
echo ""
read -p "  Are you sure you want to proceed? [y/N] " -n 1 -r < /dev/tty
echo ""
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    info "Uninstallation aborted."
    exit 0
fi

if [ ! -d "$INSTALL_DIR/stow" ]; then
    fail "Term1zz installation directory not found at $INSTALL_DIR"
fi

info "Removing Term1zz configurations..."
cd "$INSTALL_DIR/stow"
for pkg in */; do
    pkg="${pkg%/}"
    info "  Unstowing $pkg..."
    stow -D -t "$HOME" "$pkg" 2>/dev/null || true
done
ok "Term1zz configurations removed."
echo ""

info "Looking for backups..."
BACKUP_BASE="$HOME/.config/term1zz_backups"
if [ -d "$BACKUP_BASE" ]; then
    LATEST_BACKUP=$(ls -td "$BACKUP_BASE"/backup_* 2>/dev/null | head -n 1)
    if [ -n "$LATEST_BACKUP" ]; then
        info "Restoring from $LATEST_BACKUP..."
        cp -r "$LATEST_BACKUP"/* "$HOME/.config/" 2>/dev/null || true
        ok "Backup restored."
    else
        info "No backups found."
    fi
else
    info "No backup directory found."
fi
echo ""

ok "Uninstallation complete. You can now safely delete $INSTALL_DIR if desired."
echo ""
