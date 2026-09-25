# ─────────────────────────────────────────────────────────────
# Term1zz — Zsh Environment
# ─────────────────────────────────────────────────────────────

# PATH configuration
typeset -U path
path=(
    "$HOME/.local/bin"
    "/usr/local/bin"
    $path
)
export PATH

# Default Editor (respect user's existing choice, fallback to micro/nvim/nano)
if [[ -z "$EDITOR" ]]; then
    if command -v micro >/dev/null 2>&1; then
        export EDITOR="micro"
        export VISUAL="micro"
    elif command -v nvim >/dev/null 2>&1; then
        export EDITOR="nvim"
        export VISUAL="nvim"
    elif command -v nano >/dev/null 2>&1; then
        export EDITOR="nano"
        export VISUAL="nano"
    fi
fi

export LESS="-R --use-color"
if command -v bat >/dev/null 2>&1; then
    export MANPAGER="sh -c 'col -bx | bat -l man -p'"
fi

# History configuration
HISTFILE="${HOME}/.zsh_history"
HISTSIZE=50000
SAVEHIST=50000
setopt EXTENDED_HISTORY
setopt SHARE_HISTORY
setopt HIST_EXPIRE_DUPS_FIRST
setopt HIST_IGNORE_DUPS
setopt HIST_IGNORE_SPACE
setopt HIST_VERIFY

# Enable completion system
autoload -Uz compinit
if [[ -n "${ZDOTDIR:-$HOME}/.zcompdump(#qN.mh+24)" ]]; then
    compinit
else
    compinit -C
fi

# Keybindings: standard navigation
bindkey -e
bindkey '^[[H' beginning-of-line
bindkey '^[[F' end-of-line
bindkey '^[[3~' delete-char
bindkey '^[[1;5C' forward-word
bindkey '^[[1;5D' backward-word

# Interactive shell setups
if [[ -o interactive ]]; then

    # Display system overview on interactive shell launch
    if command -v fastfetch >/dev/null 2>&1; then
        fastfetch
    fi

    # Starship prompt
    if command -v starship >/dev/null 2>&1; then
        eval "$(starship init zsh)"
    fi

    # Zoxide (smarter cd)
    if command -v zoxide >/dev/null 2>&1; then
        eval "$(zoxide init --cmd cd zsh)"
    fi

    # FZF integration
    if command -v fzf >/dev/null 2>&1; then
        eval "$(fzf --zsh 2>/dev/null || true)"
    fi

    # Atuin history search
    if command -v atuin >/dev/null 2>&1; then
        eval "$(atuin init zsh --disable-up-arrow 2>/dev/null || true)"
    fi

    # ─────────────────────────────────────────────────────────
    # Modern CLI Aliases
    # ─────────────────────────────────────────────────────────

    if command -v eza >/dev/null 2>&1; then
        alias ls='eza --icons --group-directories-first'
        alias ll='eza -lah --icons --group-directories-first --git'
        alias la='eza -a --icons --group-directories-first'
        alias lt='eza --tree --level=2 --icons'
        alias ltt='eza --tree --level=3 --icons'
    fi

    if command -v bat >/dev/null 2>&1; then
        alias cat='bat --paging=never'
        alias less='bat --paging=always'
    fi

    if command -v dust >/dev/null 2>&1; then
        alias du='dust'
    fi

    if command -v btop >/dev/null 2>&1; then
        alias top='btop'
    fi

    if command -v fd >/dev/null 2>&1; then
        alias find='fd'
    fi

    if command -v rg >/dev/null 2>&1; then
        alias grep='rg'
    fi

    # Navigation shortcuts
    alias ..='cd ..'
    alias ...='cd ../..'
    alias ....='cd ../../..'
    alias .....='cd ../../../..'

    # Git shortcuts
    alias g='git'
    alias gs='git status -sb'
    alias ga='git add'
    alias gaa='git add -A'
    alias gc='git commit -m'
    alias gca='git commit --amend --no-edit'
    alias gco='git checkout'
    alias gcb='git checkout -b'
    alias gp='git push'
    alias gpf='git push --force-with-lease'
    alias gpu='git push -u origin HEAD'
    alias pl='git pull'
    alias gpl='git pull --rebase'
    alias gf='git fetch --all --prune'
    alias glog='git log --oneline --graph --decorate --all'
    alias gd='git diff'
    alias gds='git diff --staged'
    alias gb='git branch'
    alias gba='git branch -a'
    alias gbd='git branch -d'
    alias gst='git stash'
    alias gstp='git stash pop'
    alias grb='git rebase'
    alias grbi='git rebase -i'

    if command -v gitui >/dev/null 2>&1; then
        alias gg='gitui'
    elif command -v lazygit >/dev/null 2>&1; then
        alias gg='lazygit'
    fi

    # Quick shortcuts
    alias cl='clear'
    alias reload='source ~/.zshrc'
    alias ports='ss -tulnp'
    alias myip='curl -s ifconfig.me'

    # ─────────────────────────────────────────────────────────
    # Helper Functions
    # ─────────────────────────────────────────────────────────

    mkcd() {
        mkdir -p "$1" && cd "$1"
    }

    up() {
        local n=${1:-1}
        local path=""
        for ((i=1; i<=n; i++)); do
            path="../$path"
        done
        cd "$path"
    }

    gitroot() {
        local root
        root="$(git rev-parse --show-toplevel 2>/dev/null)"
        if [[ -n "$root" ]]; then
            cd "$root"
        else
            echo "Not inside a git repository."
        fi
    }

    extract() {
        if [[ -f "$1" ]]; then
            case "$1" in
                *.tar.bz2)   tar xjf "$1"     ;;
                *.tar.gz)    tar xzf "$1"     ;;
                *.tar.xz)    tar xJf "$1"     ;;
                *.tar.zst)   tar --zstd -xf "$1" ;;
                *.bz2)       bunzip2 "$1"     ;;
                *.gz)        gunzip "$1"      ;;
                *.tar)       tar xf "$1"      ;;
                *.tbz2)      tar xjf "$1"     ;;
                *.tgz)       tar xzf "$1"     ;;
                *.zip)       unzip "$1"       ;;
                *.7z)        7z x "$1"        ;;
                *.rar)       unrar x "$1"     ;;
                *.xz)        xz -d "$1"       ;;
                *.zst)       zstd -d "$1"     ;;
                *)           echo "'$1' is not recognized as a supported archive format" ;;
            esac
        else
            echo "'$1' is not a valid file"
        fi
    }
fi
