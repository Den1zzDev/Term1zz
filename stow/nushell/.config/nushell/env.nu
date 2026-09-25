# ─────────────────────────────────────────────────────────────
# Term1zz — Nushell Environment
# ─────────────────────────────────────────────────────────────

# Setup environment conversions
$env.ENV_CONVERSIONS = {
    "PATH": {
        from_string: { |s| $s | split row (char esep) | path expand --no-symlink }
        to_string: { |v| $v | path expand --no-symlink | str join (char esep) }
    }
}

# Directories to search for scripts when calling source or use
$env.NU_LIB_DIRS = [
    ($nu.default-config-dir | path join 'scripts')
]

# Directories to search for plugin binaries
$env.NU_PLUGIN_DIRS = [
    ($nu.default-config-dir | path join 'plugins')
]

# Prepend user local binaries to PATH
$env.PATH = ($env.PATH | split row (char esep) | prepend $"($env.HOME)/.local/bin" | prepend "/usr/local/bin" | uniq)

# Default editor configuration
if not ($env.EDITOR? | is-not-empty) {
    if (which micro | is-not-empty) {
        $env.EDITOR = "micro"
        $env.VISUAL = "micro"
    } else if (which nvim | is-not-empty) {
        $env.EDITOR = "nvim"
        $env.VISUAL = "nvim"
    } else if (which nano | is-not-empty) {
        $env.EDITOR = "nano"
        $env.VISUAL = "nano"
    }
}

# Starship prompt init caching
if (which starship | is-not-empty) {
    let starship_cache = ($env.HOME | path join ".cache/starship")
    mkdir $starship_cache
    starship init nu | save -f ($starship_cache | path join "init.nu")
}

# Zoxide directory jumper caching
if (which zoxide | is-not-empty) {
    let zoxide_cache = ($env.HOME | path join ".cache/zoxide")
    mkdir $zoxide_cache
    zoxide init nushell | save -f ($zoxide_cache | path join "init.nu")
}
