# ─────────────────────────────────────────────────────────────
# Term1zz — Nushell Configuration
# ─────────────────────────────────────────────────────────────

$env.config = {
    show_banner: false
    table: {
        mode: rounded
        index_mode: always
        show_empty: true
        padding: { left: 1, right: 1 }
        trim: {
            wrapping_try_keep_words: true
        }
    }
    history: {
        max_size: 50_000
        sync_on_enter: true
        file_format: "plaintext"
    }
    completions: {
        case_sensitive: false
        quick: true
        partial: true
        algorithm: "prefix"
    }
}

# Display system overview on interactive startup
if (which fastfetch | is-not-empty) {
    fastfetch
}

# Starship prompt integration
let starship_init = ($env.HOME | path join ".cache/starship/init.nu")
if ($starship_init | path exists) {
    use ($starship_init)
}

# Zoxide integration
let zoxide_init = ($env.HOME | path join ".cache/zoxide/init.nu")
if ($zoxide_init | path exists) {
    source ($zoxide_init)
}

# Modern CLI tool aliases
if (which eza | is-not-empty) {
    alias ls = eza --icons --group-directories-first
    alias ll = eza -lah --icons --group-directories-first --git
    alias la = eza -a --icons --group-directories-first
    alias lt = eza --tree --level=2 --icons
    alias ltt = eza --tree --level=3 --icons
}

if (which bat | is-not-empty) {
    alias cat = bat --paging=never
}

if (which btop | is-not-empty) {
    alias top = btop
}

if (which dust | is-not-empty) {
    alias du = dust
}

# Git shortcuts
alias g = git
alias gs = git status -sb
alias ga = git add
alias gaa = git add -A
alias gc = git commit -m
alias gp = git push
alias gpl = git pull --rebase
alias gd = git diff
alias gb = git branch
alias gco = git checkout

if (which gitui | is-not-empty) {
    alias gg = gitui
} else if (which lazygit | is-not-empty) {
    alias gg = lazygit
}

# Navigation helpers
def --env mkcd [dir: string] {
    mkdir $dir
    cd $dir
}

def --env gitroot [] {
    let root = (git rev-parse --show-toplevel | complete)
    if $root.exit_code == 0 {
        cd ($root.stdout | str trim)
    } else {
        print "Not inside a git repository."
    }
}
