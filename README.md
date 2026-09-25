# Term1zz

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
| **Catppuccin** (`catppuccin`) | Ghostty (`Catppuccin Mocha`), Fish theme, Starship palette, Zellij layout, Micro colorscheme, Bat (`Catppuccin Mocha`), Fastfetch |
| **Everforest** (`everforest`) | Ghostty (`Everforest Dark Hard`), Fish theme, Starship palette, Zellij layout, Micro (`solarized-dark`), Bat (`gruvbox-dark`), Fastfetch |
| **Tokyo Night** (`tokyonight`) | Ghostty (`TokyoNight`), Fish theme, Starship palette, Zellij layout, Micro (`tokyonight`), Bat (`TokyoNight`), Fastfetch |

---

## Included tools

### Shells and prompts
- **Fish.** Interactive shell with autosuggestions and syntax highlighting.
- **Starship.** Cross-shell prompt styled with active theme palette.
- **Fisher plugins.** Automated bootstrap of Fisher with `fzf.fish`, `autopair`, and `z`.
- **Zsh.** Configured shell with Starship and modern aliases.
- **Nushell.** Shell oriented around structured data tables.

### Multiplexers
- **Zellij.** Terminal workspace with tabs, panes, and floating layouts.
- **tmux.** Terminal multiplexer with session persistence.

### Command-line utilities
- **eza.** File listing with colors, git status, and icons.
- **bat.** Syntax-highlighted file viewer with automatic paging.
- **fd.** Fast file search utility.
- **ripgrep.** Recursive regex text search.
- **zoxide.** Directory jumper tracking frequent paths.
- **fzf.** General-purpose fuzzy finder.
- **dust.** Disk usage visualizer.
- **btop.** Terminal system resource monitor.
- **yazi.** Terminal file manager with async I/O.

### Editors and git
- **Micro.** Terminal editor with intuitive keybindings and mouse support.
- **GitUI.** Fast git terminal interface written in Rust.
- **lazygit.** Simple git terminal interface written in Go.

### Terminal emulators and fonts
- **Ghostty.** GPU-accelerated terminal emulator.
- **JetBrainsMono Nerd Font.** Automated installer for developer fonts and powerline glyphs.

### Theming
- **Catppuccin.** Soothing pastel theme preset tailored across all terminal tools.
- **Everforest.** Natural warm green theme preset tailored across all terminal tools.
- **Tokyo Night.** Clean dark neon theme preset tailored across all terminal tools.

### Productivity and search
- **fastfetch.** Maintained system information display styled with active theme accents.
- **navi.** Interactive cheatsheet browser using fzf.
- **atuin.** Searchable SQLite shell history sync.
- **xh.** Fast HTTP client with concise syntax.

---

## Dotfiles

Configurations live in `stow/` and mirror the user home directory:

```
stow/
├── fastfetch/         → ~/.config/fastfetch/config.jsonc
├── fish/              → ~/.config/fish/{config.fish, functions/, completions/}
├── ghostty/           → ~/.config/ghostty/config
├── micro/             → ~/.config/micro/settings.json
├── navi/              → ~/.config/navi/den1zz.cheat
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
