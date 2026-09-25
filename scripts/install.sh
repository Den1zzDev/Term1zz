#!/bin/sh
set -eu

# ─────────────────────────────────────────────────────────────
# Term1zz — POSIX Installer & Bootstrap
# https://github.com/Den1zzDev/Term1zz
# ─────────────────────────────────────────────────────────────

REPO="Den1zzDev/Term1zz"
INSTALL_DIR="${HOME}/.local/share/term1zz"
BIN_DIR="${HOME}/.local/bin"
TARGET_BIN="${BIN_DIR}/term1zz"

# Text styling
BOLD="\033[1m"
MAUVE="\033[38;2;203;166;247m"
GREEN="\033[38;2;166;227;161m"
RED="\033[38;2;243;139;168m"
PEACH="\033[38;2;250;179;135m"
BLUE="\033[38;2;137;180;250m"
DIM="\033[2m"
RESET="\033[0m"

info() {
    printf "%b  ▸%b %s\n" "${BLUE}${BOLD}" "${RESET}" "$1"
}

ok() {
    printf "%b  ✓%b %s\n" "${GREEN}${BOLD}" "${RESET}" "$1"
}

warn() {
    printf "%b  ⚠%b %s\n" "${PEACH}${BOLD}" "${RESET}" "$1"
}

fail() {
    printf "%b  ✗%b %s\n" "${RED}${BOLD}" "${RESET}" "$1" >&2
    exit 1
}

banner() {
    printf "\n"
    printf "%b" "${MAUVE}${BOLD}"
    printf "  +----------------------------------------+\n"
    printf "  |               Term1zz                  |\n"
    printf "  |       Modular Terminal Toolbox         |\n"
    printf "  +----------------------------------------+\n"
    printf "%b" "${RESET}"
    printf "%b  https://github.com/%s%b\n\n" "${DIM}" "${REPO}" "${RESET}"
}

detect_arch() {
    raw_arch="$(uname -m)"
    case "${raw_arch}" in
        x86_64|amd64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        armv7l|armv6l)
            echo "armv7"
            ;;
        *)
            fail "Unsupported architecture: ${raw_arch}"
            ;;
    esac
}

detect_os() {
    raw_os="$(uname -s)"
    case "${raw_os}" in
        Linux)
            echo "linux"
            ;;
        Darwin)
            echo "darwin"
            ;;
        *)
            fail "Unsupported operating system: ${raw_os}"
            ;;
    esac
}

fetch_url() {
    url="$1"
    dest="$2"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "${url}" -o "${dest}"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "${dest}" "${url}"
    else
        fail "Neither curl nor wget found. Please install one to proceed."
    fi
}

install_go_toolchain() {
    info "Go compiler not found. Installing via system package manager..."

    ESCALATE=""
    if [ "$(id -u)" -ne 0 ]; then
        if command -v sudo >/dev/null 2>&1; then
            ESCALATE="sudo"
        elif command -v run0 >/dev/null 2>&1; then
            ESCALATE="run0"
        elif command -v doas >/dev/null 2>&1; then
            ESCALATE="doas"
        else
            fail "Cannot install Go: administrative escalation tool (sudo, run0, doas) not found."
        fi
    fi

    if command -v pacman >/dev/null 2>&1; then
        ${ESCALATE} pacman -S --needed --noconfirm go
    elif command -v dnf >/dev/null 2>&1; then
        ${ESCALATE} dnf install -y golang
    elif command -v apt-get >/dev/null 2>&1; then
        ${ESCALATE} apt-get update && ${ESCALATE} apt-get install -y golang-go
    elif command -v apk >/dev/null 2>&1; then
        ${ESCALATE} apk add go
    elif command -v moss >/dev/null 2>&1; then
        ${ESCALATE} moss install -y golang
    elif command -v zypper >/dev/null 2>&1; then
        ${ESCALATE} zypper --non-interactive install go
    elif command -v brew >/dev/null 2>&1; then
        brew install go
    elif [ "${OS}" = "darwin" ]; then
        fail "Homebrew not found. Please install Homebrew from https://brew.sh first to build Term1zz on macOS."
    else
        fail "Could not find a supported package manager to install Go."
    fi
    ok "Go compiler installed."
}

banner

OS="$(detect_os)"
ARCH="$(detect_arch)"
info "Detected platform: ${OS}/${ARCH}"

# Ensure macOS standard Homebrew paths are in PATH
if [ "${OS}" = "darwin" ]; then
    for brew_path in /opt/homebrew/bin /usr/local/bin; do
        if [ -d "${brew_path}" ]; then
            case ":${PATH}:" in
                *:"${brew_path}":*) ;;
                *) export PATH="${brew_path}:${PATH}" ;;
            esac
        fi
    done
fi

# Determine whether running from inside a local checkout
SCRIPT_DIR="$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)"
LOCAL_STOW=""

if [ -d "${SCRIPT_DIR}/../stow" ]; then
    LOCAL_STOW="${SCRIPT_DIR}/../stow"
elif [ -d "${SCRIPT_DIR}/stow" ]; then
    LOCAL_STOW="${SCRIPT_DIR}/stow"
fi

# Ensure bin directory exists and is in PATH for this process
mkdir -p "${BIN_DIR}"
case ":${PATH}:" in
    *:"${BIN_DIR}":*) ;;
    *)
        export PATH="${BIN_DIR}:${PATH}"
        warn "${BIN_DIR} is not in your default PATH. Added temporarily."
        ;;
esac

# 1. Use existing binary if already compiled in local repository
if [ -n "${LOCAL_STOW}" ] && [ -x "${SCRIPT_DIR}/../term1zz" ]; then
    info "Using existing local Term1zz binary..."
    cp -f "${SCRIPT_DIR}/../term1zz" "${TARGET_BIN}"
    chmod +x "${TARGET_BIN}"
    ok "Installed Term1zz to ${TARGET_BIN}"

# 2. Compile directly if Go is already installed and local source exists
elif [ -n "${LOCAL_STOW}" ] && [ -f "${SCRIPT_DIR}/../cmd/term1zz/main.go" ] && command -v go >/dev/null 2>&1; then
    info "Building Term1zz from local source..."
    (cd "${SCRIPT_DIR}/.." && go build -o "${TARGET_BIN}" ./cmd/term1zz) || fail "Local build failed."
    ok "Built ${TARGET_BIN}"

# 3. Otherwise try to fetch official pre-built binary release from GitHub
else
    RELEASE_URL="https://github.com/${REPO}/releases/latest/download/term1zz-${OS}-${ARCH}.tar.gz"
    TMP_DIR="$(mktemp -d 2>/dev/null || mktemp -d -t 'term1zz')"
    cleanup() {
        rm -rf "${TMP_DIR}"
    }
    trap cleanup EXIT INT TERM

    ARCHIVE="${TMP_DIR}/term1zz.tar.gz"

    info "Fetching latest Term1zz release..."
    if fetch_url "${RELEASE_URL}" "${ARCHIVE}" 2>/dev/null; then
        tar -xzf "${ARCHIVE}" -C "${TMP_DIR}" || fail "Failed to extract release archive."
        if [ -f "${TMP_DIR}/term1zz" ]; then
            mv -f "${TMP_DIR}/term1zz" "${TARGET_BIN}"
            chmod +x "${TARGET_BIN}"
            ok "Installed Term1zz to ${TARGET_BIN}"
        else
            fail "Release archive did not contain term1zz binary."
        fi
    else
        # 4. GitHub release not available (e.g. testing in fresh container before first GH tag)
        info "Pre-built release not found on GitHub."

        if ! command -v go >/dev/null 2>&1; then
            install_go_toolchain
        fi

        if [ -n "${LOCAL_STOW}" ] && [ -f "${SCRIPT_DIR}/../cmd/term1zz/main.go" ]; then
            info "Building Term1zz from local source..."
            (cd "${SCRIPT_DIR}/.." && go build -o "${TARGET_BIN}" ./cmd/term1zz) || fail "Build failed."
            ok "Built ${TARGET_BIN}"
        else
            info "Installing Term1zz via go install..."
            GOBIN="${BIN_DIR}" go install "github.com/${REPO}/cmd/term1zz@latest" || fail "Go install failed."
            ok "Installed Term1zz via go install."
        fi
    fi
fi

# Synchronize dotfiles repository
mkdir -p "${INSTALL_DIR}"
if [ -n "${LOCAL_STOW}" ]; then
    info "Using local configurations from ${LOCAL_STOW}."
    mkdir -p "${INSTALL_DIR}/stow"
    cp -R "${LOCAL_STOW}"/* "${INSTALL_DIR}/stow/" 2>/dev/null || true
else
    info "Updating configuration templates at ${INSTALL_DIR}..."
    CONFIGS_TAR="https://github.com/${REPO}/archive/refs/heads/main.tar.gz"
    TMP_SRC="${TMP_DIR}/source.tar.gz"
    if fetch_url "${CONFIGS_TAR}" "${TMP_SRC}" 2>/dev/null; then
        mkdir -p "${TMP_DIR}/src"
        tar -xzf "${TMP_SRC}" -C "${TMP_DIR}/src" --strip-components=1
        if [ -d "${TMP_DIR}/src/stow" ]; then
            mkdir -p "${INSTALL_DIR}/stow"
            cp -R "${TMP_DIR}/src/stow"/* "${INSTALL_DIR}/stow/"
            ok "Configurations synchronized."
        fi
    fi
fi

ok "Term1zz is ready."
printf "\n"

# Run Term1zz
# In non-interactive contexts (e.g. redirected or without terminal), run directly.
# When interactive, attach to TTY if needed.
if [ -t 0 ]; then
    exec "${TARGET_BIN}" --stow-dir "${INSTALL_DIR}/stow" "$@"
elif (exec < /dev/tty) 2>/dev/null; then
    exec "${TARGET_BIN}" --stow-dir "${INSTALL_DIR}/stow" "$@" < /dev/tty
else
    exec "${TARGET_BIN}" --stow-dir "${INSTALL_DIR}/stow" "$@"
fi
