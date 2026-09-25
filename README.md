# Term1zz (Reborn)

Terminal environment manager and dotfile orchestrator for Linux and macOS.

Term1zz detects your operating system, package manager, and privilege escalation tools to install modern command-line utilities, apply unified theme presets, and link dotfiles.

---

## Quickstart

### One-line installer

```sh
curl -fsSL https://raw.githubusercontent.com/Den1zzDev/Term1zz/main/scripts/install.sh | sh
```

The installer downloads a pre-compiled binary when available. If no binary matches your system, it compiles from source automatically.

### Build from source

```sh
# Build binary
go build -o term1zz ./cmd/term1zz

# Run unit tests
go test -v ./...

# Preview planned actions without modifying system state
./term1zz --dry-run

# Run interactive TUI
./term1zz

# Run batch mode with a specific theme preset
./term1zz --batch --theme everforest
```

---

## Keybindings

| Key | Action |
| --- | --- |
| `↑` / `k`, `↓` / `j` | Move selection cursor |
| `Space` | Toggle item (exclusive selection in `[6] Theming`) |
| `m` | Cycle install mode (`[all]` tool + config, `[pkg]` package only, `[cfg]` config only) |
| `a` | Toggle all items in category (selects highlighted theme in `[6] Theming`) |
| `1` - `7` | Jump directly to category |
| `Tab` / `Shift+Tab` | Cycle through categories |
| `d` | Toggle dry-run mode |
| `Enter` | Review installation plan |
| `y` (in review) | Confirm and execute plan |
| `q` / `Ctrl+C` | Quit |

---

## Suite theming

When you select a theme preset in `[6] Theming` or pass `--theme <name>`, Term1zz configures the entire terminal suite specifically for that palette:

| Preset | Target tools tailored |
| --- | --- |
| **Catppuccin** (`catppuccin`) | Ghostty (`Catppuccin Mocha`), Fish theme, Starship palette, Zellij layout, Micro (`catppuccin-mocha`), Bat (`Catppuccin Mocha`), Fastfetch |
| **Everforest** (`everforest`) | Ghostty (`Everforest Dark Hard`), Fish theme, Starship palette, Zellij layout, Micro (`solarized-dark`), Bat (`gruvbox-dark`), Fastfetch |
| **Tokyo Night** (`tokyonight`) | Ghostty (`TokyoNight`), Fish theme, Starship palette, Zellij layout, Micro (`tokyonight`), Bat (`TwoDark`), Fastfetch |

---

## Included tools

### Shells and prompts
- **Fish.** Interactive shell with autosuggestions and syntax highlighting.
- **Starship.** Cross-shell prompt styled with active theme palette.
- **Fisher plugins.** Automated bootstrap of Fisher with `fzf.fish` and `autopair`.
- **Zsh.** Configured shell with Starship, modern history deduplication, and completion caching.
- **Nushell.** Shell oriented around structured data tables with Starship prompt.

### Multiplexers
- **Zellij.** Terminal workspace with tabs, panes, and floating layouts.
- **tmux.** Terminal multiplexer with session persistence.

### Command-line utilities
- **eza.** Modern `ls` with colors, git status, and tree view.
- **bat.** Syntax-highlighted file viewer with automatic paging and git integration.
- **fd.** Fast, user-friendly alternative to `find`.
- **ripgrep.** Recursive regex text search respecting `.gitignore`.
- **zoxide.** Smarter `cd` tracking frequently visited paths.
- **fzf.** General-purpose fuzzy finder.
- **dust.** Graphical disk usage visualizer.
- **btop.** Terminal system resource monitor.
- **yazi.** Terminal file manager with async I/O and image previews.

### Editors and git
- **Micro.** Terminal editor with intuitive keybindings and mouse support.
- **GitUI.** Fast git terminal interface written in Rust.
- **lazygit.** Terminal UI for git commands with visual branch management.

### Terminal emulators and fonts
- **Ghostty.** GPU-accelerated terminal emulator with native tabs and font features.
- **JetBrainsMono Nerd Font.** Automated installer for developer fonts and powerline glyphs with native package mapping and direct archive fallback.

### Theming
- **Catppuccin.** Soothing pastel theme preset tailored across all terminal tools.
- **Everforest.** Natural warm green theme preset tailored across all terminal tools.
- **Tokyo Night.** Clean dark neon theme preset tailored across all terminal tools.

### Productivity and search
- **fastfetch.** Maintained system information display styled with active theme accents.
- **navi.** Interactive cheatsheet browser using fzf.
- **atuin.** Searchable SQLite shell history sync.
- **xh.** User-friendly HTTP client with concise syntax and colors.

---

## Supported platforms and package managers

Term1zz supports both Linux and macOS with automatic detection of packages and non-destructive privilege escalation:

| Platform / Distro | Package Manager | Escalation |
| --- | --- | --- |
| **macOS** | Homebrew (`brew`) | None (Homebrew standard) |
| **AerynOS** | `moss` | `run0` / `sudo` |
| **Arch Linux** | `pacman` | `sudo` / `run0` / `doas` |
| **Fedora** | `dnf` | `sudo` / `run0` / `doas` |
| **Ubuntu / Debian** | `apt-get` | `sudo` / `run0` / `doas` |
| **Alpine Linux** | `apk` | `sudo` / `doas` |
| **openSUSE** | `zypper` | `sudo` / `run0` / `doas` |
| **Void Linux** | `xbps-install` | `sudo` / `doas` |

---

## Dotfiles

Configurations live in `stow/` and mirror the user home directory:

```
stow/
├── fastfetch/         → ~/.config/fastfetch/config.jsonc
├── fish/              → ~/.config/fish/{config.fish, conf.d/, functions/, completions/}
├── ghostty/           → ~/.config/ghostty/config
├── micro/             → ~/.config/micro/settings.json
├── navi/              → ~/.config/navi/den1zz.cheat
├── nushell/           → ~/.config/nushell/{config.nu, env.nu}
├── starship/          → ~/.config/starship.toml
├── theme-catppuccin/  → ~/.config/{ghostty,starship.toml,zellij,micro,bat,fastfetch,fish}
├── theme-everforest/  → ~/.config/{ghostty,starship.toml,zellij,micro,bat,fastfetch,fish}
├── theme-tokyonight/  → ~/.config/{ghostty,starship.toml,zellij,micro,bat,fastfetch,fish}
├── zellij/            → ~/.config/zellij/config.kdl
└── zsh/               → ~/.zshrc
```

Term1zz creates symlinks from `$HOME` to these files. If a file already exists at the destination, Term1zz moves it to a timestamped directory in `~/.local/state/term1zz/backups/`.

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development workflows, package mapping guidelines, and container test commands.

---

## License

[MIT](LICENSE)
