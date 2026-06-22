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

INTERACTIVE=0
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -i|--interactive) INTERACTIVE=1 ;;
        -h|--help)
            echo "Usage: setup.sh [OPTIONS]"
            echo "Options:"
            echo "  -i, --interactive    Prompt before stowing each configuration"
            echo "  -h, --help           Show this help message"
            exit 0
            ;;
        *) echo "Unknown parameter: $1"; exit 1 ;;
    esac
    shift
done

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

# ── Warning & Verification ──────────────────────
echo -e "${RED}${BOLD}  ⚠ WARNING: Destructive Operation (Backups Taken)${RESET}"
echo -e "  This installer will overwrite existing terminal configurations"
echo -e "  (Fish, Ghostty, Starship, Zellij, Micro, etc.) that conflict with Term1zz."
echo -e "  A backup of your conflicting files will be automatically created under"
echo -e "  ~/.config/term1zz_backups/ before they are replaced."
echo ""

# Read from /dev/tty to allow interactive prompt when piped via curl | bash
read -p "  Are you sure you want to proceed? [y/N] " -n 1 -r < /dev/tty
echo ""
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    info "Installation aborted by user. No files were modified."
    exit 0
fi

# ── Custom Installers ───────────────────────────
install_nerd_font() {
    if ! fc-list | grep -qi "Nerd Font"; then
        info "No Nerd Font detected. Installing JetBrainsMono Nerd Font..."
        mkdir -p ~/.local/share/fonts
        wget -qO /tmp/JetBrainsMono.zip "https://github.com/ryanoasis/nerd-fonts/releases/latest/download/JetBrainsMono.zip"
        unzip -qo /tmp/JetBrainsMono.zip -d ~/.local/share/fonts/
        fc-cache -f
        rm -f /tmp/JetBrainsMono.zip
        ok "JetBrainsMono Nerd Font installed."
    else
        ok "Nerd Font already installed."
    fi
}

# ── Preflight ───────────────────────────────────
info "Detecting privilege escalation tool..."
if command -v run0 &>/dev/null; then
    ELEVATE_CMD="run0 --background="
elif command -v doas &>/dev/null; then
    ELEVATE_CMD="doas"
elif command -v sudo &>/dev/null; then
    ELEVATE_CMD="sudo"
else
    if [ "$EUID" -eq 0 ]; then
        ELEVATE_CMD=""
    else
        fail "No privilege escalation tool found (run0, doas, or sudo)."
    fi
fi
if [ -n "$ELEVATE_CMD" ]; then
    ok "Using ${ELEVATE_CMD%% *} for privilege escalation."
else
    ok "Running as root."
fi
echo ""

info "Detecting package manager..."
if command -v pacman &>/dev/null; then
    if command -v yay &>/dev/null; then
        PKG_MGR="yay"
    elif command -v paru &>/dev/null; then
        PKG_MGR="paru"
    else
        PKG_MGR="pacman"
    fi
elif command -v dnf &>/dev/null; then
    PKG_MGR="dnf"
elif command -v zypper &>/dev/null; then
    PKG_MGR="zypper"
elif command -v apt-get &>/dev/null; then
    PKG_MGR="apt-get"
else
    fail "No supported package manager found. Supported: pacman, yay, paru, dnf, zypper, apt-get."
fi
ok "Detected $PKG_MGR."
echo ""

# ── Step 1: Install packages ───────────────────
info "Installing packages via $PKG_MGR..."

# Universal base
PACKAGES=(git stow fish micro curl unzip fontconfig fzf zoxide ripgrep btop mpv)
if [ "$PKG_MGR" != "apt-get" ]; then
    PACKAGES+=(bat fd dust)
else
    PACKAGES+=(fd-find)
fi

if [[ "$PKG_MGR" == "pacman" || "$PKG_MGR" == "yay" || "$PKG_MGR" == "paru" ]]; then
    PACKAGES+=(eza fastfetch zellij atuin gitui xh)
    echo -e "  ${DIM}${PACKAGES[*]}${RESET}"
    if [ "$PKG_MGR" == "pacman" ]; then
        $ELEVATE_CMD pacman -S --needed --noconfirm "${PACKAGES[@]}" || fail "Package installation failed."
    else
        $PKG_MGR -S --needed --noconfirm "${PACKAGES[@]}" || fail "Package installation failed."
    fi

elif [ "$PKG_MGR" == "dnf" ]; then
    # Fedora: explicitly configure Terra repository first
    $ELEVATE_CMD dnf install -y dnf-plugins-core
    $ELEVATE_CMD dnf config-manager --add-repo https://terra.fyralabs.com/terra.repo
    
    PACKAGES+=(eza fastfetch zellij atuin gitui)
    echo -e "  ${DIM}${PACKAGES[*]}${RESET}"
    $ELEVATE_CMD dnf install -y "${PACKAGES[@]}" || fail "Package installation failed."

elif [ "$PKG_MGR" == "zypper" ]; then
    PACKAGES+=(eza fastfetch zellij)
    echo -e "  ${DIM}${PACKAGES[*]}${RESET}"
    $ELEVATE_CMD zypper install -y "${PACKAGES[@]}" || fail "Package installation failed."

elif [ "$PKG_MGR" == "apt-get" ]; then
    $ELEVATE_CMD apt-get update
    
    # apt-get needs specific dependencies for adding repos
    APT_BASE=(bat gpg wget software-properties-common curl "${PACKAGES[@]}")
    echo -e "  ${DIM}${APT_BASE[*]}${RESET}"
    $ELEVATE_CMD apt-get install -y "${APT_BASE[@]}" || fail "Base package installation failed."
    
    # Resolve bat binary conflict
    mkdir -p ~/.local/bin && ln -sf /usr/bin/batcat ~/.local/bin/bat
    mkdir -p ~/.local/bin && ln -sf /usr/bin/fdfind ~/.local/bin/fd
    
    # Add third-party repositories for eza and fastfetch
    $ELEVATE_CMD mkdir -p /etc/apt/keyrings
    wget -qO- https://raw.githubusercontent.com/eza-community/eza/main/deb.asc | $ELEVATE_CMD gpg --dearmor --yes -o /etc/apt/keyrings/gierens.gpg
    echo "deb [signed-by=/etc/apt/keyrings/gierens.gpg] http://deb.gierens.de stable main" | $ELEVATE_CMD tee /etc/apt/sources.list.d/gierens.list > /dev/null
    $ELEVATE_CMD chmod 644 /etc/apt/keyrings/gierens.gpg /etc/apt/sources.list.d/gierens.list
    
    $ELEVATE_CMD add-apt-repository -y ppa:zhangsongcui3371/fastfetch
    
    $ELEVATE_CMD apt-get update
    $ELEVATE_CMD apt-get install -y eza fastfetch || fail "eza and fastfetch installation failed."
    
    info "Installing zellij directly from official binary..."
    wget -qO- "https://github.com/zellij-org/zellij/releases/latest/download/zellij-x86_64-unknown-linux-musl.tar.gz" | tar -xz -C ~/.local/bin zellij || fail "Zellij download failed."
fi

ok "All packages installed."
echo ""

# ── Custom Tool Installations ────────────────────
info "Installing font and missing tools..."
install_nerd_font

if ! command -v atuin &>/dev/null; then
    curl --proto '=https' --tlsv1.2 -LsSf https://setup.atuin.sh | sh
fi
if ! command -v xh &>/dev/null; then
    curl -sfL https://raw.githubusercontent.com/ducaale/xh/master/install.sh | sh
fi
if ! command -v zed &>/dev/null; then
    curl -f https://zed.dev/install.sh | sh
fi
echo ""

# ── Standalone Installations ───────────────────
info "Checking for starship prompt..."
if ! command -v starship &>/dev/null; then
    curl -sS https://starship.rs/install.sh | sh -s -- -y || fail "Starship installation failed."
    ok "Starship installed."
else
    ok "Starship is already installed."
fi
echo ""

# ── Step 2: Clone or update repository ─────────
if [ -d "$INSTALL_DIR/.git" ]; then
    info "Repository already exists at ${DIM}$INSTALL_DIR${RESET}"
    info "Pulling latest changes..."
    git -C "$INSTALL_DIR" reset --hard >/dev/null 2>&1
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

info "Cleaning up legacy configurations..."
find "$HOME/.config" -type l -lname '*Den1zzfiles*' -delete 2>/dev/null || true
find "$HOME" -maxdepth 1 -type l -lname '*Den1zzfiles*' -delete 2>/dev/null || true

# Backup existing configs
BACKUP_DIR="$HOME/.config/term1zz_backups/backup_$(date +%s)"
info "Backing up conflicting configurations to ${DIM}$BACKUP_DIR${RESET}..."
mkdir -p "$BACKUP_DIR"
for pkg in */; do
    pkg="${pkg%/}"
    if [ -d "$HOME/.config/$pkg" ] && [ ! -L "$HOME/.config/$pkg" ]; then
        cp -r "$HOME/.config/$pkg" "$BACKUP_DIR/"
    elif [ -f "$HOME/.config/$pkg" ] && [ ! -L "$HOME/.config/$pkg" ]; then
        cp "$HOME/.config/$pkg" "$BACKUP_DIR/"
    fi
    if [ -f "$HOME/.$pkg" ] && [ ! -L "$HOME/.$pkg" ]; then
        cp "$HOME/.$pkg" "$BACKUP_DIR/"
    fi
done
ok "Backup created."

# Stow each package individually for clearer error reporting
for pkg in */; do
    pkg="${pkg%/}"  # Remove trailing slash
    
    if [ "$INTERACTIVE" -eq 1 ]; then
        read -p "  Do you want to stow $pkg? [Y/n] " -n 1 -r </dev/tty
        echo ""
        if [[ $REPLY =~ ^[Nn]$ ]]; then
            info "  Skipping $pkg."
            echo ""
            continue
        fi
    fi

    info "  Stowing ${MAUVE}${pkg}${RESET}..."
    stow --adopt -t "$HOME" -R "$pkg" 2>&1 || warn "Failed to stow $pkg — conflicts may exist."
    echo ""
done

# Reset the git repository to discard adopted user files and enforce the Term1zz versions
git -C "$INSTALL_DIR" reset --hard >/dev/null 2>&1

ok "All selected configurations linked."
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
            echo "$FISH_PATH" | $ELEVATE_CMD tee -a /etc/shells >/dev/null
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
