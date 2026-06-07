#!/bin/bash
set -e

# ═══════════════════════════════════════════════
#   Term1zz — Remote Installer
#   https://codeberg.org/Den1zz/Term1zz
#
#   Usage:
#     curl -sL https://codeberg.org/Den1zz/Term1zz/raw/branch/main/setup.sh | bash
# ═══════════════════════════════════════════════

REPO_URL="https://codeberg.org/Den1zz/Term1zz.git"
INSTALL_DIR="$HOME/.local/share/Term1zz"

# ── Colors ──────────────────────────────────────
BOLD='\033[1m'
DIM='\033[2m'
MAUVE='\033[38;2;203;166;247m'
GREEN='\033[38;2;166;227;161m'
RED='\033[38;2;243;139;168m'
PEACH='\033[38;2;250;179;135m'
BLUE='\033[38;2;137;180;250m'
RESET='\033[0m'

info()  { echo -e "${BLUE}${BOLD}  ▸${RESET} $1"; }
ok()    { echo -e "${GREEN}${BOLD}  ✓${RESET} $1"; }
warn()  { echo -e "${PEACH}${BOLD}  ⚠${RESET} $1"; }
fail()  { echo -e "${RED}${BOLD}  ✗${RESET} $1"; exit 1; }

# ── Banner ──────────────────────────────────────
echo ""
echo -e "${MAUVE}${BOLD}"
echo "  ╔════════════════════════════════════════╗"
echo "  ║                                        ║"
echo "  ║       ✦  T e r m 1 z z  ✦             ║"
echo "  ║       Modular Terminal Framework        ║"
echo "  ║                                        ║"
echo "  ╚════════════════════════════════════════╝"
echo -e "${RESET}"
echo -e "  ${DIM}https://codeberg.org/Den1zz/Term1zz${RESET}"
echo ""

# ── Preflight ───────────────────────────────────
if ! command -v pacman &>/dev/null; then
    fail "pacman not found. This installer requires an Arch-based distribution."
fi

# ── Step 1: Install packages ───────────────────
PACKAGES=(
    git
    stow
    fish
    zellij
    micro
    eza
    bat
    fastfetch
    starship
)

info "Installing packages via pacman..."
echo -e "  ${DIM}${PACKAGES[*]}${RESET}"
echo ""

sudo pacman -S --needed --noconfirm "${PACKAGES[@]}" || fail "Package installation failed."
ok "All packages installed."
echo ""

# ── Step 2: Clone or update repository ─────────
if [ -d "$INSTALL_DIR/.git" ]; then
    info "Repository already exists at ${DIM}$INSTALL_DIR${RESET}"
    info "Pulling latest changes..."
    git -C "$INSTALL_DIR" pull --rebase || fail "Git pull failed."
    ok "Repository updated."
else
    info "Cloning Term1zz to ${DIM}$INSTALL_DIR${RESET}..."
    mkdir -p "$(dirname "$INSTALL_DIR")"
    git clone "$REPO_URL" "$INSTALL_DIR" || fail "Git clone failed."
    ok "Repository cloned."
fi
echo ""

# ── Step 3: Stow all packages ─────────────────
info "Linking configurations with GNU Stow..."

if [ ! -d "$INSTALL_DIR/stow" ]; then
    fail "stow/ directory not found in repository. Is the repo structure correct?"
fi

cd "$INSTALL_DIR/stow"

# Ensure target directories exist
mkdir -p "$HOME/.config"

# Stow each package individually for clearer error reporting
for pkg in */; do
    pkg="${pkg%/}"  # Remove trailing slash
    info "  Stowing ${MAUVE}${pkg}${RESET}..."
    stow -t "$HOME" -R "$pkg" 2>&1 || warn "Failed to stow $pkg — conflicts may exist."
done

ok "All configurations linked."
echo ""

# ── Step 3b: Bootstrap Fish plugins via Fisher ──
info "Bootstrapping Fish plugins (Fisher)..."
if command -v fish &>/dev/null; then
    fish -c "curl -sL https://raw.githubusercontent.com/jorgebucaran/fisher/main/functions/fisher.fish | source && fisher update" \
        || warn "Fisher bootstrap failed — run 'fisher update' manually inside fish."
    ok "Fish plugins installed."
else
    warn "fish not found in PATH — skipping plugin bootstrap."
fi
echo ""

# ── Step 4: Set default shell ──────────────────
FISH_PATH="$(which fish 2>/dev/null)"

if [ -n "$FISH_PATH" ]; then
    CURRENT_SHELL="$(basename "$SHELL")"
    if [ "$CURRENT_SHELL" != "fish" ]; then
        info "Setting fish as the default shell..."
        # Ensure fish is in /etc/shells
        if ! grep -qxF "$FISH_PATH" /etc/shells 2>/dev/null; then
            echo "$FISH_PATH" | sudo tee -a /etc/shells >/dev/null
        fi
        chsh -s "$FISH_PATH" || warn "Could not change default shell. Run: chsh -s $FISH_PATH"
        ok "Default shell set to fish."
    else
        ok "Fish is already the default shell."
    fi
else
    warn "Fish shell not found in PATH. Skipping shell change."
fi
echo ""

# ── Done ────────────────────────────────────────
echo -e "${GREEN}${BOLD}"
echo "  ╔════════════════════════════════════════╗"
echo "  ║                                        ║"
echo "  ║       ✦  Setup Complete!  ✦            ║"
echo "  ║                                        ║"
echo "  ╚════════════════════════════════════════╝"
echo -e "${RESET}"
echo -e "  ${DIM}Configs installed to: $INSTALL_DIR${RESET}"
echo -e "  ${DIM}Symlinks managed by: GNU Stow${RESET}"
echo ""
echo -e "  ${MAUVE}${BOLD}→${RESET} Restart your terminal or log out and back in"
echo -e "    for all changes to take effect."
echo ""
