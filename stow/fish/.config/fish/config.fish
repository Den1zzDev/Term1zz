# ─────────────────────────────────────────────────────────────
# Term1zz — Fish Shell Environment
# ─────────────────────────────────────────────────────────────

# Ensure user binaries are in PATH
fish_add_path "$HOME/.local/bin"
fish_add_path "/usr/local/bin"

# Default Editor (respect user's existing choice, fallback to micro/nvim/nano)
if not set -q EDITOR
    if command -q micro
        set -gx EDITOR micro
        set -gx VISUAL micro
    else if command -q nvim
        set -gx EDITOR nvim
        set -gx VISUAL nvim
    else if command -q nano
        set -gx EDITOR nano
        set -gx VISUAL nano
    end
end

set -gx NAVI_PATH "$HOME/.config/navi"
set -gx LESS "-R --use-color"

if command -q bat
    set -gx MANPAGER "sh -c 'col -bx | bat -l man -p'"
end

# Fuzzy Finder (fzf) Defaults
if command -q fd
    set -gx FZF_DEFAULT_COMMAND "fd --type f --hidden --follow --exclude .git"
    set -gx FZF_CTRL_T_COMMAND "$FZF_DEFAULT_COMMAND"
    set -gx FZF_ALT_C_COMMAND "fd --type d --hidden --follow --exclude .git"
end

set -gx FZF_DEFAULT_OPTS " \
  --height=50% --layout=reverse --border=rounded --info=inline \
  --prompt='  ' --pointer='▶' --marker='✓' \
  --preview-window='right:50%:wrap'"

# ─────────────────────────────────────────────────────────────
# Interactive Session Initialization
# ─────────────────────────────────────────────────────────────
if status is-interactive

    # Disable default greeting
    set -g fish_greeting ""

    # Display system overview on interactive shell launch
    if command -q fastfetch
        fastfetch
    end

    # Starship Prompt
    if command -q starship
        starship init fish | source
    end

    # Zoxide (smarter directory jumping)
    if command -q zoxide
        zoxide init --cmd cd fish | source
    end

    # Fuzzy Finder keybindings (Ctrl-R history, Ctrl-T files, Alt-C directories)
    if command -q fzf
        fzf --fish 2>/dev/null | source
    end

    # Navi cheatsheets (Ctrl-G)
    if command -q navi
        navi widget fish 2>/dev/null | source
    end

    # Atuin history search
    if command -q atuin
        atuin init fish --disable-up-arrow 2>/dev/null | source
    end

    # ─────────────────────────────────────────────────────────
    # Modern CLI Aliases
    # ─────────────────────────────────────────────────────────

    if command -q eza
        alias ls  "eza --icons --group-directories-first"
        alias ll  "eza -lah --icons --group-directories-first --git"
        alias la  "eza -a --icons --group-directories-first"
        alias lt  "eza --tree --level=2 --icons"
        alias ltt "eza --tree --level=3 --icons"
    end

    if command -q bat
        alias cat  "bat --paging=never"
        alias less "bat --paging=always"
    end

    if command -q dust
        alias du "dust"
    end

    if command -q btop
        alias top "btop"
    end

    if command -q fd
        alias find "fd"
    end

    if command -q rg
        alias grep "rg"
    end

    # Navigation shortcuts
    alias ..    "cd .."
    alias ...   "cd ../.."
    alias ....  "cd ../../.."
    alias ..... "cd ../../../.."

    # ─────────────────────────────────────────────────────────
    # Abbreviations (auto-expand on Space)
    # ─────────────────────────────────────────────────────────

    # Git
    abbr g      "git"
    abbr gs     "git status -sb"
    abbr ga     "git add"
    abbr gaa    "git add -A"
    abbr gc     "git commit -m"
    abbr gca    "git commit --amend --no-edit"
    abbr gco    "git checkout"
    abbr gcb    "git checkout -b"
    abbr gp     "git push"
    abbr gpf    "git push --force-with-lease"
    abbr gpu    "git push -u origin HEAD"
    abbr pl     "git pull"
    abbr gpl    "git pull --rebase"
    abbr gf     "git fetch --all --prune"
    abbr glog   "git log --oneline --graph --decorate --all"
    abbr gd     "git diff"
    abbr gds    "git diff --staged"
    abbr gb     "git branch"
    abbr gba    "git branch -a"
    abbr gbd    "git branch -d"
    abbr gst    "git stash"
    abbr gstp   "git stash pop"
    abbr grb    "git rebase"
    abbr grbi   "git rebase -i"

    if command -q gitui
        abbr gg "gitui"
    else if command -q lazygit
        abbr gg "lazygit"
    end

    # GitHub CLI
    abbr ghpr   "gh pr create"
    abbr ghprl  "gh pr list"
    abbr ghprv  "gh pr view --web"

    # Shortcuts
    abbr cl     "clear"
    abbr reload "source ~/.config/fish/config.fish"
    abbr ports  "ss -tulnp"
    abbr myip   "curl -s ifconfig.me"

    # ─────────────────────────────────────────────────────────
    # Utility Functions
    # ─────────────────────────────────────────────────────────

    # Create directory and jump into it
    function mkcd --description "Create a directory and cd into it"
        mkdir -p $argv[1] && cd $argv[1]
    end

    # Move up N directories
    function up --description "Navigate up N directories"
        set -l n 1
        if test (count $argv) -gt 0
            set n $argv[1]
        end
        set -l path ""
        for i in (seq 1 $n)
            set path "../$path"
        end
        cd $path
    end

    # Jump to git repository root
    function gitroot --description "Navigate to git repository root"
        set -l root (git rev-parse --show-toplevel 2>/dev/null)
        if test -n "$root"
            cd "$root"
        else
            echo "Not inside a git repository."
        end
    end

    # Universal archive extractor
    function extract --description "Extract any common archive format"
        if test -f "$argv[1]"
            switch "$argv[1]"
                case "*.tar.bz2"
                    tar xjf "$argv[1]"
                case "*.tar.gz"
                    tar xzf "$argv[1]"
                case "*.tar.xz"
                    tar xJf "$argv[1]"
                case "*.tar.zst"
                    tar --zstd -xf "$argv[1]"
                case "*.bz2"
                    bunzip2 "$argv[1]"
                case "*.gz"
                    gunzip "$argv[1]"
                case "*.tar"
                    tar xf "$argv[1]"
                case "*.tbz2"
                    tar xjf "$argv[1]"
                case "*.tgz"
                    tar xzf "$argv[1]"
                case "*.zip"
                    unzip "$argv[1]"
                case "*.7z"
                    7z x "$argv[1]"
                case "*.rar"
                    unrar x "$argv[1]"
                case "*.xz"
                    xz -d "$argv[1]"
                case "*.zst"
                    zstd -d "$argv[1]"
                case "*"
                    echo "'$argv[1]' is not recognized as a supported archive format"
            end
        else
            echo "'$argv[1]' is not a valid file"
        end
    end

end
