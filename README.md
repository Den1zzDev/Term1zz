<div align="center">

# ✦ Den1zzfiles ✦

*A curated collection of personal dotfiles & system configurations*

[![KDE](https://img.shields.io/badge/DE-KDE_Plasma-5c91cf?style=flat-square&logo=kde&logoColor=white)](https://kde.org)
[![Shell](https://img.shields.io/badge/Shell-Fish_%2F_Oksh-89b4fa?style=flat-square&logo=gnu-bash&logoColor=white)](https://fishshell.com)
[![Ghostty](https://img.shields.io/badge/Terminal-Ghostty-a6e3a1?style=flat-square)](https://ghostty.org)
[![Catppuccin](https://img.shields.io/badge/Theme-Catppuccin-f5c2e7?style=flat-square)](https://catppuccin.com)

</div>

---

## 🚀 Installation

> [!WARNING]
> The included `install.sh` script is **EXPERIMENTAL** and highly opinionated. It is designed for fresh CachyOS / KDE installations and will aggressively modify configurations in your home directory. **It is not recommended to run this on an existing setup.** Please review the script manually before execution!

```bash
cd ~/.config
git clone https://codeberg.org/Den1zz/Den1zzfiles
cd Den1zzfiles
./install.sh
```

---

## 🛠️ Included Configurations

| Component | Choice | Notes |
|-----------|--------|-------|
| **Desktop Environment** | [KDE Plasma](https://kde.org) | Wayland |
| **Shell** | Fish / Oksh | Atuin · zoxide · fzf integrations |
| **Prompt** | Starship | Catppuccin styling |
| **Terminal** | Ghostty | — |
| **Editor** | Zed / Micro | — |
| **UI / Theming** | [Catppuccin KDE](https://github.com/catppuccin/kde) | KDE Plasma theming |

---

## 🐚 Shell Comparison — Fish vs Oksh

Both shells are configured and ready to use. Pick whichever fits your workflow.

| | [Fish](https://fishshell.com) | [Oksh](https://github.com/ibara/oksh) |
|---|---|---|
| **Philosophy** | Friendly & feature-rich out of the box | Minimal & POSIX-compliant |
| **Syntax** | Custom (not POSIX) | POSIX sh / ksh |
| **Autosuggestions** | ✅ Built-in | ❌ |
| **Syntax highlighting** | ✅ Built-in | ❌ |
| **Tab completions** | ✅ Extensive | Basic |
| **Scripting** | Own syntax — not portable | POSIX — portable everywhere |
| **Startup speed** | Fast | Faster |
| **Resource usage** | Moderate | Very light |
| **Best for** | Daily interactive use | Lightweight sessions, scripting, POSIX compatibility |

---

## 💛 Credits & Acknowledgements

These configs are built on the shoulders of a bunch of great open-source projects. Full credit to their authors.

### 🖥️ Environment

| Project | What it does |
|---------|-------------|
| [KDE Plasma](https://kde.org) | Powerful and customizable desktop environment |
| [Fish Shell](https://fishshell.com) / [Oksh](https://github.com/ibara/oksh) | Fish for interactive use, Oksh for lightweight POSIX compliance |
| [Starship](https://starship.rs) | Cross-shell prompt |
| [Ghostty](https://ghostty.org) | Fast, native terminal emulator |
| [Zed](https://zed.dev) | High-performance code editor |
| [Micro](https://micro-editor.github.io) | Modern terminal text editor |
| [Catppuccin KDE](https://github.com/catppuccin/kde) | Soothing pastel theme for KDE Plasma |
| [Catppuccin](https://catppuccin.com) | Soothing pastel color scheme (used everywhere) |

### 🔧 CLI Toolchain

| Project | Replaces | What it does |
|---------|----------|-------------|
| [eza](https://github.com/eza-community/eza) | `ls` | Modern file lister with icons & git info |
| [bat](https://github.com/sharkdp/bat) | `cat` / `less` | Syntax-highlighted file viewer |
| [fd](https://github.com/sharkdp/fd) | `find` | Fast, user-friendly file finder |
| [ripgrep](https://github.com/BurntSushi/ripgrep) | `grep` | Blazing-fast recursive search |
| [dust](https://github.com/bootandy/dust) | `du` | Intuitive disk usage viewer |
| [btop](https://github.com/aristocratos/btop) | `top` | Resource monitor with a beautiful TUI |
| [zoxide](https://github.com/ajeetdsouza/zoxide) | `cd` | Smarter directory jumping |
| [fzf](https://github.com/junegunn/fzf) | — | General-purpose fuzzy finder |
| [Atuin](https://github.com/atuinsh/atuin) | shell history | Synced, searchable shell history |
| [uutils coreutils](https://github.com/uutils/coreutils) | GNU coreutils | Dynamically maps `uu-*` binaries — falls back to system coreutils if not installed |
| [bfetch](https://github.com/Mjoyufull/bfetch) | [fastfetch](https://github.com/fastfetch-cli/fastfetch) | System fetch on shell startup — falls back to fastfetch if not installed |

### 🎥 Media

| Project | What it does |
|---------|-------------|
| [streamlink](https://streamlink.github.io) | Extracts streams from sites like Twitch |
| [mpv](https://mpv.io) | Minimal, powerful media player |

<div align="center">

*Feel free to steal, fork, or adapt anything here.*

</div>
