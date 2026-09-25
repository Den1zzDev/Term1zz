# ─────────────────────────────────────────────────────────────────────────────
# Term1zz — Zsh Environment
# ─────────────────────────────────────────────────────────────────────────────

# History configuration
HISTFILE="${HOME}/.zsh_history"
HISTSIZE=10000
SAVEHIST=10000
setopt APPEND_HISTORY
setopt SHARE_HISTORY
setopt HIST_IGNORE_DUPS
setopt HIST_IGNORE_SPACE

# Prompt initialization
if command -v starship >/dev/null 2>&1; then
    eval "$(starship init zsh)"
fi

# Directory jumping
if command -v zoxide >/dev/null 2>&1; then
    eval "$(zoxide init zsh)"
fi

# Modern CLI Aliases
if command -v eza >/dev/null 2>&1; then
    alias ls='eza --icons'
    alias ll='eza -la --icons --git'
    alias lt='eza --tree --level=2 --icons'
fi

if command -v bat >/dev/null 2>&1; then
    alias cat='bat --paging=never'
fi

if command -v fzf >/dev/null 2>&1; then
    eval "$(fzf --zsh 2>/dev/null || true)"
fi
