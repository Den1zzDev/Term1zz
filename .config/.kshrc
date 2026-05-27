# ~/.kshrc — oksh interactive shell configuration
# Sourced automatically for every interactive oksh session via $ENV.
# ──────────────────────────────────────────────

# ──────────────────────────────────────────────
# HISTORY SETTINGS
# ──────────────────────────────────────────────
export HISTFILE="$HOME/.oksh_history"
export HISTSIZE=10000
export HISTCONTROL="ignoredups:ignorespace"

# ──────────────────────────────────────────────
# ENVIRONMENT VARIABLES
# ──────────────────────────────────────────────
export VISUAL="zeditor"
export EDITOR="zeditor"
export TERM="xterm-256color"
export NAVI_PATH="$HOME/.config/navi"

# bat → Catppuccin Mocha theme everywhere
export BAT_THEME="Catppuccin Mocha"

# fzf → Catppuccin Mocha palette + bat preview
export FZF_DEFAULT_OPTS="--height=50% --layout=reverse --border=rounded --info=inline \
--prompt='  ' --pointer='▶' --marker='✓' \
--color=bg+:#313244,bg:#1e1e2e,spinner:#f5e0dc,hl:#f38ba8 \
--color=fg:#cdd6f4,header:#f38ba8,info:#cba6f7,pointer:#f5e0dc \
--color=marker:#f5e0dc,fg+:#cdd6f4,prompt:#cba6f7,hl+:#f38ba8 \
--preview='bat --style=numbers --color=always {}' \
--preview-window='right:50%:wrap'"
export FZF_DEFAULT_COMMAND="fd --type f --hidden --follow --exclude .git"
export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"
export FZF_ALT_C_COMMAND="fd --type d --hidden --follow --exclude .git"

# less → colours
export LESS="-R --use-color"
export MANPAGER="sh -c 'col -bx | bat -l man -p'"

# ──────────────────────────────────────────────
# PATH
# ──────────────────────────────────────────────
case ":$PATH:" in
    *:"$HOME/.local/bin":*) ;;
    *) export PATH="$HOME/.local/bin:$PATH" ;;
esac

# ──────────────────────────────────────────────
# GREETING
# ──────────────────────────────────────────────
clear
if command -v fastfetch >/dev/null 2>&1; then
    fastfetch
elif command -v bfetch >/dev/null 2>&1; then
    bfetch
fi

# ──────────────────────────────────────────────
# PROMPT: Starship
#
# `starship init bash` generates bash-specific code (PROMPT_COMMAND,
# DEBUG trap, arrays) that oksh cannot eval. A PS1 subshell captures
# $? then calls starship prompt directly.
# ──────────────────────────────────────────────
if command -v starship >/dev/null 2>&1; then
    # POSIX subshells inherit $? from the parent at the moment of
    # evaluation, so $? here is the user's last exit code. starship
    # prompt outputs the fully-rendered prompt string to stdout, which
    # becomes the value of PS1 for that prompt cycle.
    PS1='$(starship prompt --status=$?)'
fi

# ──────────────────────────────────────────────
# SMART CD: zoxide
# ──────────────────────────────────────────────
if command -v zoxide >/dev/null 2>&1; then
    eval "$(zoxide init posix --cmd cd --hook prompt)"
fi

# ══════════════════════════════════════════════
# ALIASES — Modern CLI replacements
# ══════════════════════════════════════════════

# ls → eza
if command -v eza >/dev/null 2>&1; then
    alias ls='eza --icons --color=always --group-directories-first'
    alias ll='eza -lah --icons --color=always --group-directories-first --git'
    alias la='eza -a --icons --color=always --group-directories-first'
    alias lt='eza --tree --level=2 --icons --color=always'
    alias ltt='eza --tree --level=3 --icons --color=always'
    alias lg='eza -lah --icons --color=always --git --git-repos'
fi

# cat → bat
if command -v bat >/dev/null 2>&1; then
    alias cat='bat'
    alias less='bat --paging=always'
fi

# du → dust
if command -v dust >/dev/null 2>&1; then
    alias du='dust'
fi

# top → btop
if command -v btop >/dev/null 2>&1; then
    alias top='btop'
fi

# find → fd
if command -v fd >/dev/null 2>&1; then
    alias find='fd'
fi

# grep → ripgrep
if command -v rg >/dev/null 2>&1; then
    alias grep='rg'
fi

# nano → micro
if command -v micro >/dev/null 2>&1; then
    alias nano='micro'
fi

# ══════════════════════════════════════════════
# uutils coreutils — Dynamic Mapping
#
# Aliases every /usr/bin/uu-* binary, skipping shell builtins and
# the exclude list (tools already aliased to modern replacements).
# ══════════════════════════════════════════════
_uutils_exclude="ls cat du top find grep"
_uutils_builtins="echo test true false pwd printf kill type"

for _uu_bin in /usr/bin/uu-*; do
    [ -x "$_uu_bin" ] || continue

    _base="${_uu_bin##*/}"   # uu-cp
    _base="${_base#uu-}"     # cp

    # Skip POSIX builtins
    _skip=0
    for _b in $_uutils_builtins; do
        if [ "$_base" = "$_b" ]; then _skip=1; break; fi
    done
    [ "$_skip" -eq 1 ] && continue

    # Skip commands already aliased to modern tools
    _skip=0
    for _x in $_uutils_exclude; do
        if [ "$_base" = "$_x" ]; then _skip=1; break; fi
    done
    [ "$_skip" -eq 1 ] && continue

    # Apply special safety flags
    case "$_base" in
        cp|mv|rm) alias "$_base"="$_uu_bin -i"  ;;
        mkdir)    alias "$_base"="$_uu_bin -pv"  ;;
        *)        alias "$_base"="$_uu_bin"       ;;
    esac
done

unset _uu_bin _base _skip _b _x _uutils_exclude _uutils_builtins

# ══════════════════════════════════════════════
# sudo → run0 wrapper
#
# Routes sudo calls through run0. When the target command has a
# uu-* uutils equivalent, that binary is used instead.
# ══════════════════════════════════════════════
sudo() {
    _s_exclude="ls cat du top find grep"

    if [ $# -eq 0 ]; then
        command run0 --background=""
        return
    fi

    _s_cmd="$1"

    # Forward flag-only invocations straight to run0
    case "$_s_cmd" in
        -*)
            command run0 --background="" "$@"
            return
            ;;
    esac

    shift  # remove cmd from "$@", rest are its args

    # Check for a uutils binary, respecting the exclude list
    if [ -x "/usr/bin/uu-$_s_cmd" ]; then
        _s_excl=0
        for _s_x in $_s_exclude; do
            if [ "$_s_cmd" = "$_s_x" ]; then _s_excl=1; break; fi
        done

        if [ "$_s_excl" -eq 0 ]; then
            case "$_s_cmd" in
                cp|mv|rm)
                    command run0 --background="" "/usr/bin/uu-$_s_cmd" -i "$@"
                    return ;;
                mkdir)
                    command run0 --background="" "/usr/bin/uu-$_s_cmd" -pv "$@"
                    return ;;
                *)
                    command run0 --background="" "/usr/bin/uu-$_s_cmd" "$@"
                    return ;;
            esac
        fi
    fi

    command run0 --background="" "$_s_cmd" "$@"
}

# ══════════════════════════════════════════════
# NAVIGATION
# ══════════════════════════════════════════════
alias ..='cd ..'
alias ...='cd ../..'
alias ....='cd ../../..'
alias .....='cd ../../../..'
alias cdh='cd ~'
alias cdb='cd -'

# ══════════════════════════════════════════════
# GIT ALIASES
# ══════════════════════════════════════════════
alias g='git'
alias gg='gitui'
alias gi='git init'
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
alias gcp='git cherry-pick'
alias gtag='git tag'
alias greset='git reset --hard HEAD'
alias gclean='git clean -fd'

# — GitHub CLI —
alias ghpr='gh pr create'
alias ghprl='gh pr list'
alias ghprv='gh pr view --web'

# ══════════════════════════════════════════════
# PACKAGE MANAGER ALIASES
# ══════════════════════════════════════════════

# npm
alias ni='npm install'
alias nid='npm install --save-dev'
alias nig='npm install -g'
alias nr='npm run'
alias nrd='npm run dev'
alias nrb='npm run build'
alias nrt='npm run test'
alias nrc='npm run check'
alias nrp='npm run preview'

# pnpm
alias pi='pnpm install'
alias pa='pnpm add'
alias pad='pnpm add -D'
alias pr='pnpm run'
alias prd='pnpm run dev'
alias prb='pnpm run build'

# yarn
alias yi='yarn install'
alias ya='yarn add'
alias yr='yarn run'

# ══════════════════════════════════════════════
# PYTHON ALIASES
# ══════════════════════════════════════════════
alias py='python3'
alias pip='pip3'
alias venv='python3 -m venv .venv && . .venv/bin/activate'
alias va='. .venv/bin/activate'
alias vd='deactivate'

# ══════════════════════════════════════════════
# DOCKER ALIASES
# ══════════════════════════════════════════════
alias d='docker'
alias dc='docker compose'
alias dcu='docker compose up -d'
alias dcd='docker compose down'
alias dcl='docker compose logs -f'
alias dps='docker ps'
alias dpsa='docker ps -a'
alias dimg='docker images'
alias dex='docker exec -it'
alias drm='docker rm'
alias drmi='docker rmi'
alias dprune='docker system prune -af --volumes'

# ══════════════════════════════════════════════
# KUBERNETES ALIASES
# ══════════════════════════════════════════════
alias k='kubectl'
alias kgp='kubectl get pods'
alias kgs='kubectl get svc'
alias kgn='kubectl get nodes'
alias kd='kubectl describe'
alias kl='kubectl logs -f'
alias kex='kubectl exec -it'
alias kctx='kubectx'
alias kns='kubens'
alias kaf='kubectl apply -f'
alias kdf='kubectl delete -f'

# ══════════════════════════════════════════════
# MISCELLANEOUS ALIASES
# ══════════════════════════════════════════════
alias e='"$EDITOR" .'
alias http='xh'
alias nv='navi'
alias q='exit'
alias cl='clear'
alias reload='. ~/.kshrc'
alias kshrc='"$EDITOR" ~/.kshrc'
alias path='echo "$PATH" | tr ":" "\n"'
alias ports='ss -tulnp'
alias myip='curl -s ifconfig.me'
alias weather='curl -s wttr.in'
alias speedtest='curl -s https://raw.githubusercontent.com/sivel/speedtest-cli/master/speedtest.py | python3 -'

# ══════════════════════════════════════════════
# FUNCTIONS
# ══════════════════════════════════════════════

# mkdir + cd
mkcd() {
    mkdir -p "$1" && cd "$1"
}

# go up N directories
up() {
    _up_n="${1:-1}"
    _up_path=""
    _up_i=0
    while [ "$_up_i" -lt "$_up_n" ]; do
        _up_path="../$_up_path"
        _up_i=$(( _up_i + 1 ))
    done
    cd "$_up_path"
}

# universal archive extractor
extract() {
    if [ -f "$1" ]; then
        case "$1" in
            *.tar.bz2)  tar xjf "$1"            ;;
            *.tar.gz)   tar xzf "$1"             ;;
            *.tar.xz)   tar xJf "$1"             ;;
            *.tar.zst)  tar --zstd -xf "$1"      ;;
            *.bz2)      bunzip2 "$1"             ;;
            *.gz)       gunzip "$1"              ;;
            *.tar)      tar xf "$1"              ;;
            *.tbz2)     tar xjf "$1"             ;;
            *.tgz)      tar xzf "$1"             ;;
            *.zip)      unzip "$1"               ;;
            *.7z)       7z x "$1"                ;;
            *.rar)      unrar x "$1"             ;;
            *.xz)       xz -d "$1"               ;;
            *.zst)      zstd -d "$1"             ;;
            *)          echo "'$1' — unknown archive format" ;;
        esac
    else
        echo "'$1' is not a valid file"
    fi
}

# fzf-powered directory jump (zoxide fallback)
fcd() {
    _fcd_dir="$(fd --type d --hidden --exclude .git | fzf +m)"
    [ -n "$_fcd_dir" ] && cd "$_fcd_dir"
}

# interactively kill a process via fzf
fkill() {
    _fkill_pid="$(ps -eo pid,comm,args | fzf --header="Select process to kill" | awk '{print $1}')"
    if [ -n "$_fkill_pid" ]; then
        kill -9 "$_fkill_pid" && echo "Killed PID $_fkill_pid"
    fi
}

# cd to git repo root
gitroot() {
    cd "$(git rev-parse --show-toplevel)"
}

# pick a known SSH host via fzf
sshfzf() {
    _sfzf_host="$(grep -E "^Host " "$HOME/.ssh/config" 2>/dev/null | awk '{print $2}' | fzf --prompt=" SSH host: ")"
    [ -n "$_sfzf_host" ] && ssh "$_sfzf_host"
}

# watch Twitch stream via streamlink + mpv
twitch() {
    streamlink --player "mpv" "https://www.twitch.tv/$1" best
}

# search and watch anime via ani-cli
anime() {
    ani-cli "$@"
}

# ══════════════════════════════════════════════
# SSH AGENT AUTO-START
# ══════════════════════════════════════════════
if ! pgrep -xu "$USER" ssh-agent >/dev/null 2>&1; then
    eval "$(ssh-agent -s)"
    export SSH_AUTH_SOCK
    export SSH_AGENT_PID
fi
