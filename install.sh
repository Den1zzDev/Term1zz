#!/bin/bash
set -e

echo "================================================="
echo "   ✦ Den1zzfiles Installer for CachyOS / KDE ✦"
echo "================================================="
echo "   ⚠️  WARNING: THIS SCRIPT IS EXPERIMENTAL! ⚠️"
echo "   Please review the code before running it."
echo "   It will modify configurations in your home"
echo "   directory and install various packages."
echo "================================================="
echo ""
read -p "Are you sure you want to proceed? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborting installation."
    exit 1
fi
echo ""

# Ensure we are in the correct directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_CONFIG_DIR="$SCRIPT_DIR/.config"

if [ ! -d "$REPO_CONFIG_DIR" ]; then
    echo "Error: .config directory not found in the current folder."
    echo "Please run this script from the root of the Den1zzfiles repository."
    exit 1
fi

echo "[1/3] Installing required CLI tools and dependencies..."

# List of packages based on the README and config files
PACKAGES=(
    # Shells
    fish
    oksh
    
    # Prompt & Terminal
    starship
    ghostty
    
    # Editors
    zed
    micro
    
    # Core CLI Tools / Modern Replacements
    fastfetch
    zoxide
    fzf
    atuin
    eza
    bat
    dust
    btop
    fd
    ripgrep
    uutils-coreutils
    
    # Media
    streamlink
    mpv
)

# Detect AUR helper available in CachyOS (paru is default, yay is an alternative)
AUR_HELPER=""
if command -v paru >/dev/null 2>&1; then
    AUR_HELPER="paru"
elif command -v yay >/dev/null 2>&1; then
    AUR_HELPER="yay"
fi

PACMAN_PKGS=()
AUR_PKGS=()

echo "Categorizing packages for installation..."
for pkg in "${PACKAGES[@]}"; do
    # Check if pacman knows about the package in the official repos or groups
    if pacman -Si "$pkg" >/dev/null 2>&1 || pacman -Sg "$pkg" >/dev/null 2>&1; then
        PACMAN_PKGS+=("$pkg")
    else
        AUR_PKGS+=("$pkg")
    fi
done

if [ ${#PACMAN_PKGS[@]} -gt 0 ]; then
    echo "Installing official repository packages using pacman..."
    run0 pacman -S --needed --noconfirm "${PACMAN_PKGS[@]}"
fi

if [ ${#AUR_PKGS[@]} -gt 0 ]; then
    if [ -n "$AUR_HELPER" ]; then
        echo "Installing AUR packages using $AUR_HELPER..."
        $AUR_HELPER -S --needed --noconfirm "${AUR_PKGS[@]}"
    else
        echo "Warning: The following packages are not in official repos and no AUR helper (paru/yay) was found:"
        for pkg in "${AUR_PKGS[@]}"; do
            echo "  - $pkg"
        done
        echo "Please install them manually."
    fi
fi

echo ""
echo "[2/3] Setting up dotfiles..."

# Backup existing configs just in case
BACKUP_DIR="$HOME/.config/dotfiles_backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

echo "Backing up existing configurations and symlinking new ones..."

# Ensure .config exists
mkdir -p "$HOME/.config"

# Process .config items
for item in fastfetch fish ghostty starship.toml; do
    TARGET="$HOME/.config/$item"
    
    # Check if target exists or is a broken symlink
    if [ -e "$TARGET" ] || [ -L "$TARGET" ]; then
        mv "$TARGET" "$BACKUP_DIR/"
    fi
    
    # Use -snf to avoid dereferencing existing symlinks to directories
    ln -snf "$REPO_CONFIG_DIR/$item" "$TARGET"
done

# Process home directory items
for item in .kshrc .profile; do
    TARGET="$HOME/$item"
    
    # Check if target exists or is a broken symlink
    if [ -e "$TARGET" ] || [ -L "$TARGET" ]; then
        mv "$TARGET" "$BACKUP_DIR/"
    fi
    
    # Use -snf for safety
    ln -snf "$REPO_CONFIG_DIR/$item" "$TARGET"
done

echo ""
echo "[3/3] Finalizing setup..."

# Change default shell
echo ""
echo "Which shell would you like to set as your default?"
echo "1) fish (feature-rich & friendly)"
echo "2) oksh (minimal & POSIX-compliant)"
echo "3) Do not change default shell"
read -p "Select an option [1-3]: " shell_choice
echo

case $shell_choice in
    1)
        if grep -q "/usr/bin/fish" /etc/shells; then
            chsh -s /usr/bin/fish
            echo "Default shell changed to fish."
        else
            echo "Error: /usr/bin/fish not found in /etc/shells."
        fi
        ;;
    2)
        if grep -q "/usr/bin/oksh" /etc/shells; then
            chsh -s /usr/bin/oksh
            echo "Default shell changed to oksh."
        else
            echo "Error: /usr/bin/oksh not found in /etc/shells."
        fi
        ;;
    *)
        echo "Default shell unchanged."
        ;;
esac

echo ""
echo "================================================="
echo "   Setup Complete!"
echo "   Please restart your terminal or log out and"
echo "   log back in for all changes to take effect."
echo "================================================="
