<div align="center">

# ✦ Term1zz ✦

*A modular terminal framework — managed by GNU Stow*

[![Arch](https://img.shields.io/badge/Distro-Arch_Linux-1793d1?style=flat-square&logo=archlinux&logoColor=white)](https://archlinux.org)
[![Fedora](https://img.shields.io/badge/Distro-Fedora-294172?style=flat-square&logo=fedora&logoColor=white)](https://fedoraproject.org)
[![OpenSUSE](https://img.shields.io/badge/Distro-OpenSUSE-73ba25?style=flat-square&logo=opensuse&logoColor=white)](https://www.opensuse.org)
[![Debian](https://img.shields.io/badge/Distro-Debian-A81D33?style=flat-square&logo=debian&logoColor=white)](https://www.debian.org)
[![Ubuntu](https://img.shields.io/badge/Distro-Ubuntu-E95420?style=flat-square&logo=ubuntu&logoColor=white)](https://ubuntu.com)
[![Shell](https://img.shields.io/badge/Shell-Fish-89b4fa?style=flat-square&logo=gnu-bash&logoColor=white)](https://fishshell.com)
[![Ghostty](https://img.shields.io/badge/Terminal-Ghostty-a6e3a1?style=flat-square)](https://ghostty.org)
[![Zellij](https://img.shields.io/badge/Multiplexer-Zellij-fab387?style=flat-square)](https://zellij.dev)
[![Catppuccin](https://img.shields.io/badge/Theme-Catppuccin-f5c2e7?style=flat-square)](https://catppuccin.com)

</div>

---

## ⚡ One-Line Install

> [!NOTE]
> Supports **Arch, Fedora, OpenSUSE, Debian, and Ubuntu**. The script will automatically configure third-party repositories (like the Terra repository for Fedora and official PPAs for Debian/Ubuntu) to guarantee access to modern utilities. It will then install packages and symlink configurations into your home directory. Review the script before running!

```bash
curl -sL https://codeberg.org/Den1zz/Term1zz/raw/branch/main/setup.sh | bash
```

<details>
<summary>What does it do?</summary>

1. Installs dependencies via your native package manager (`pacman`, `dnf`, `zypper`, or `apt-get`) — `git`, `stow`, `fish`, `zellij`, `micro`, `eza`, `bat`, `fastfetch`, and `starship`
2. Clones the repository to `~/.local/share/Term1zz`
3. Uses GNU Stow to symlink all configuration packages into `$HOME`
4. Sets Fish as the default shell

</details>

---

## 🗂️ Repository Structure

Configurations are organized as independent **Stow packages**. Each package mirrors the home directory layout so `stow -t $HOME <pkg>` creates correct symlinks.

```
stow/
├── fastfetch/  → ~/.config/fastfetch/config.jsonc
├── fish/       → ~/.config/fish/{config.fish, themes/, functions/, ...}
├── ghostty/    → ~/.config/ghostty/config
├── micro/      → ~/.config/micro/settings.json
├── navi/       → ~/.config/navi/den1zz.cheat
├── starship/   → ~/.config/starship.toml
└── zellij/     → ~/.config/zellij/config.kdl
```

---

## 🛠️ Included Configurations

| Component | Choice | Notes |
|-----------|--------|-------|
| **Shell** | [Fish](https://fishshell.com) | Atuin · zoxide · fzf integrations |
| **Prompt** | [Starship](https://starship.rs) | Catppuccin Frappé styling |
| **Terminal** | [Ghostty](https://ghostty.org) | Fast, native emulator |
| **Multiplexer** | [Zellij](https://zellij.dev) | Catppuccin Mocha theme, Alt-key bindings |
| **Editor** | [Micro](https://micro-editor.github.io) | Terminal editor with mouse support |
| **GUI Editor** | [Zed](https://zed.dev) | High-performance code editor |
| **Theming** | [Catppuccin](https://catppuccin.com) | Soothing pastel scheme — everywhere |

---

## 💛 Credits & Acknowledgements

These configs are built on the shoulders of a bunch of great open-source projects. Full credit to their authors.

### 🖥️ Environment

| Project | What it does |
|---------|-------------|
| [Fish Shell](https://fishshell.com) | Feature-rich, friendly interactive shell |
| [Starship](https://starship.rs) | Cross-shell prompt |
| [Ghostty](https://ghostty.org) | Fast, native terminal emulator |
| [Zellij](https://zellij.dev) | Modern terminal multiplexer with a plugin system |
| [Micro](https://micro-editor.github.io) | Modern terminal text editor — intuitive & mouse-friendly |
| [Zed](https://zed.dev) | High-performance code editor |
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
| [navi](https://github.com/denisidoro/navi) | cheatsheets | Interactive cheatsheet searcher using `fzf` |
| [xh](https://github.com/ducaale/xh) | `curl` / `httpie` | Fast, user-friendly HTTP client |
| [gitui](https://github.com/extrawurst/gitui) | `lazygit` | Blazing-fast terminal UI for Git |

### 🎥 Media

| Project | What it does |
|---------|-------------|
| [streamlink](https://streamlink.github.io) | Extracts streams from sites like Twitch |
| [mpv](https://mpv.io) | Minimal, powerful media player |
| [ani-cli](https://github.com/pystardust/ani-cli) | Search and stream anime from the terminal |

<div align="center">

*Feel free to steal, fork, or adapt anything here.*

</div>
