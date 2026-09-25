package registry

import "github.com/Den1zzDev/Term1zz/internal/distro"

// AllCategories returns the ordered list of categories.
func AllCategories() []Category {
	return []Category{
		CatShells,
		CatMultiplex,
		CatCoreutils,
		CatEditorsGit,
		CatTerminals,
		CatThemes,
		CatUtils,
	}
}

// Items returns the catalog of all tools managed by Term1zz.
func Items() []Item {
	return []Item{
		// ─────────────────────────────────────────────
		// 1. Shells & Prompts
		// ─────────────────────────────────────────────
		{
			ID:          "fish",
			Name:        "Fish Shell",
			Description: "Interactive command-line shell with autosuggestions and syntax highlighting",
			Category:    CatShells,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "fish",
				distro.PMPacman: "fish",
				distro.PMDnf:    "fish",
				distro.PMApt:    "fish",
				distro.PMApk:    "fish",
				distro.PMZypper: "fish",
				distro.PMXbps:   "fish-shell",
				distro.PMBrew:   "fish",
			},
			BinaryName:       "fish",
			StowPackage:      "fish",
			DefaultSelected:  true,
			PostInstallNotes: "Pre-configured with completions, aliases, and theme integration.",
		},
		{
			ID:          "starship",
			Name:        "Starship Prompt",
			Description: "Minimal, fast, and customizable prompt for any shell",
			Category:    CatShells,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "starship",
				distro.PMPacman: "starship",
				distro.PMDnf:    "starship",
				distro.PMApt:    "starship",
				distro.PMApk:    "starship",
				distro.PMZypper: "starship",
				distro.PMXbps:   "starship",
				distro.PMBrew:   "starship",
			},
			FallbackRepo:    "starship/starship",
			BinaryName:      "starship",
			StowPackage:     "starship",
			DefaultSelected: true,
		},
		{
			ID:          "fisher",
			Name:        "Fisher Plugins (Fish)",
			Description: "Plugin manager for Fish installing fzf.fish, autopair, and tide",
			Category:    CatShells,
			CustomScript: "if command -v fish >/dev/null 2>&1; then " +
				"fish -c 'curl -sL https://raw.githubusercontent.com/jorgebucaran/fisher/main/functions/fisher.fish | source && fisher install jorgebucaran/fisher PatrickF1/fzf.fish jethrokuan/z' || true; " +
				"fi",
			DefaultSelected:  true,
			PostInstallNotes: "Bootstraps Fisher and recommended plugins directly inside Fish.",
		},
		{
			ID:          "zsh",
			Name:        "Zsh Shell",
			Description: "Powerful shell with extensive customization and plugin support",
			Category:    CatShells,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "zsh",
				distro.PMPacman: "zsh",
				distro.PMDnf:    "zsh",
				distro.PMApt:    "zsh",
				distro.PMApk:    "zsh",
				distro.PMZypper: "zsh",
				distro.PMXbps:   "zsh",
				distro.PMBrew:   "zsh",
			},
			BinaryName:      "zsh",
			StowPackage:     "zsh",
			DefaultSelected: false,
		},
		{
			ID:          "nushell",
			Name:        "Nushell",
			Description: "Modern shell that treats command output as structured data tables",
			Category:    CatShells,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "nushell",
				distro.PMPacman: "nushell",
				distro.PMDnf:    "nushell",
				distro.PMApk:    "nushell",
				distro.PMBrew:   "nushell",
			},
			FallbackRepo:    "nushell/nushell",
			BinaryName:      "nu",
			DefaultSelected: false,
		},

		// ─────────────────────────────────────────────
		// 2. Multiplexers
		// ─────────────────────────────────────────────
		{
			ID:          "zellij",
			Name:        "Zellij",
			Description: "Terminal workspace with tabs, panes, floating windows, and layouts",
			Category:    CatMultiplex,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "zellij",
				distro.PMPacman: "zellij",
				distro.PMDnf:    "zellij",
				distro.PMApk:    "zellij",
				distro.PMXbps:   "zellij",
				distro.PMBrew:   "zellij",
			},
			FallbackRepo:     "zellij-org/zellij",
			BinaryName:       "zellij",
			StowPackage:      "zellij",
			DefaultSelected:  true,
			PostInstallNotes: "Configured with Alt-key navigation and custom layouts.",
		},
		{
			ID:          "tmux",
			Name:        "tmux",
			Description: "Terminal multiplexer with session persistence and detachment",
			Category:    CatMultiplex,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "tmux",
				distro.PMPacman: "tmux",
				distro.PMDnf:    "tmux",
				distro.PMApt:    "tmux",
				distro.PMApk:    "tmux",
				distro.PMZypper: "tmux",
				distro.PMXbps:   "tmux",
				distro.PMBrew:   "tmux",
			},
			BinaryName:      "tmux",
			DefaultSelected: false,
		},

		// ─────────────────────────────────────────────
		// 3. CLI Core Replacements
		// ─────────────────────────────────────────────
		{
			ID:          "eza",
			Name:        "eza (modern ls)",
			Description: "Modern replacement for ls with colors, icons, git status, and tree view",
			Category:    CatCoreutils,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "eza",
				distro.PMPacman: "eza",
				distro.PMDnf:    "eza",
				distro.PMApt:    "eza",
				distro.PMApk:    "eza",
				distro.PMZypper: "eza",
				distro.PMXbps:   "eza",
				distro.PMBrew:   "eza",
			},
			FallbackRepo:    "eza-community/eza",
			BinaryName:      "eza",
			DefaultSelected: true,
		},
		{
			ID:          "bat",
			Name:        "bat (modern cat)",
			Description: "Cat clone with syntax highlighting, git integration, and paging",
			Category:    CatCoreutils,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "bat",
				distro.PMPacman: "bat",
				distro.PMDnf:    "bat",
				distro.PMApt:    "bat",
				distro.PMApk:    "bat",
				distro.PMZypper: "bat",
				distro.PMXbps:   "bat",
				distro.PMBrew:   "bat",
			},
			FallbackRepo:    "sharkdp/bat",
			BinaryName:      "bat",
			DefaultSelected: true,
		},
		{
			ID:          "fd",
			Name:        "fd (modern find)",
			Description: "Fast, user-friendly alternative to find",
			Category:    CatCoreutils,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "fd",
				distro.PMPacman: "fd",
				distro.PMDnf:    "fd-find",
				distro.PMApt:    "fd-find",
				distro.PMApk:    "fd",
				distro.PMZypper: "fd",
				distro.PMXbps:   "fd",
				distro.PMBrew:   "fd",
			},
			FallbackRepo:    "sharkdp/fd",
			BinaryName:      "fd",
			DefaultSelected: true,
		},
		{
			ID:          "ripgrep",
			Name:        "ripgrep (rg)",
			Description: "Recursive line search that respects gitignore by default",
			Category:    CatCoreutils,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "ripgrep",
				distro.PMPacman: "ripgrep",
				distro.PMDnf:    "ripgrep",
				distro.PMApt:    "ripgrep",
				distro.PMApk:    "ripgrep",
				distro.PMZypper: "ripgrep",
				distro.PMXbps:   "ripgrep",
				distro.PMBrew:   "ripgrep",
			},
			FallbackRepo:    "BurntSushi/ripgrep",
			BinaryName:      "rg",
			DefaultSelected: true,
		},
		{
			ID:          "zoxide",
			Name:        "zoxide (smarter cd)",
			Description: "Learns visited paths to jump directories in few keystrokes",
			Category:    CatCoreutils,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "zoxide",
				distro.PMPacman: "zoxide",
				distro.PMDnf:    "zoxide",
				distro.PMApt:    "zoxide",
				distro.PMApk:    "zoxide",
				distro.PMZypper: "zoxide",
				distro.PMXbps:   "zoxide",
				distro.PMBrew:   "zoxide",
			},
			FallbackRepo:    "ajeetdsouza/zoxide",
			BinaryName:      "zoxide",
			DefaultSelected: true,
		},
		{
			ID:          "fzf",
			Name:        "fzf",
			Description: "General-purpose command-line fuzzy finder",
			Category:    CatCoreutils,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "fzf",
				distro.PMPacman: "fzf",
				distro.PMDnf:    "fzf",
				distro.PMApt:    "fzf",
				distro.PMApk:    "fzf",
				distro.PMZypper: "fzf",
				distro.PMXbps:   "fzf",
				distro.PMBrew:   "fzf",
			},
			FallbackRepo:    "junegunn/fzf",
			BinaryName:      "fzf",
			DefaultSelected: true,
		},
		{
			ID:          "dust",
			Name:        "dust (visual du)",
			Description: "Intuitive graphical disk usage auditor in the terminal",
			Category:    CatCoreutils,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "dust",
				distro.PMPacman: "dust",
				distro.PMDnf:    "dust",
				distro.PMApk:    "dust",
				distro.PMZypper: "dust",
				distro.PMXbps:   "dust",
				distro.PMBrew:   "dust",
			},
			FallbackRepo:    "bootandy/dust",
			BinaryName:      "dust",
			DefaultSelected: true,
		},
		{
			ID:          "btop",
			Name:        "btop",
			Description: "Resource monitor displaying CPU, memory, disks, network, and processes",
			Category:    CatCoreutils,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "btop",
				distro.PMPacman: "btop",
				distro.PMDnf:    "btop",
				distro.PMApt:    "btop",
				distro.PMApk:    "btop",
				distro.PMZypper: "btop",
				distro.PMXbps:   "btop",
				distro.PMBrew:   "btop",
			},
			FallbackRepo:    "aristocratos/btop",
			BinaryName:      "btop",
			DefaultSelected: true,
		},
		{
			ID:          "yazi",
			Name:        "yazi",
			Description: "Fast terminal file manager with async I/O and image previews",
			Category:    CatCoreutils,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "yazi",
				distro.PMPacman: "yazi",
				distro.PMDnf:    "yazi",
				distro.PMApk:    "yazi",
				distro.PMBrew:   "yazi",
			},
			FallbackRepo:    "sxyazi/yazi",
			BinaryName:      "yazi",
			DefaultSelected: true,
		},

		// ─────────────────────────────────────────────
		// 4. Editors & Git
		// ─────────────────────────────────────────────
		{
			ID:          "micro",
			Name:        "Micro Editor",
			Description: "Intuitive terminal text editor with mouse support and keybindings",
			Category:    CatEditorsGit,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "micro",
				distro.PMPacman: "micro",
				distro.PMDnf:    "micro",
				distro.PMApt:    "micro",
				distro.PMApk:    "micro",
				distro.PMZypper: "micro",
				distro.PMXbps:   "micro",
				distro.PMBrew:   "micro",
			},
			FallbackRepo:    "zyedidia/micro",
			BinaryName:      "micro",
			StowPackage:     "micro",
			DefaultSelected: true,
		},
		{
			ID:          "gitui",
			Name:        "GitUI",
			Description: "Fast terminal interface for git operations written in Rust",
			Category:    CatEditorsGit,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "gitui",
				distro.PMPacman: "gitui",
				distro.PMDnf:    "gitui",
				distro.PMApk:    "gitui",
				distro.PMXbps:   "gitui",
				distro.PMBrew:   "gitui",
			},
			FallbackRepo:    "extrawurst/gitui",
			BinaryName:      "gitui",
			DefaultSelected: true,
		},
		{
			ID:          "lazygit",
			Name:        "lazygit",
			Description: "Terminal UI for git commands with visual branch management",
			Category:    CatEditorsGit,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "lazygit",
				distro.PMPacman: "lazygit",
				distro.PMDnf:    "lazygit",
				distro.PMApt:    "lazygit",
				distro.PMApk:    "lazygit",
				distro.PMZypper: "lazygit",
				distro.PMBrew:   "lazygit",
			},
			FallbackRepo:    "jesseduffield/lazygit",
			BinaryName:      "lazygit",
			DefaultSelected: false,
		},

		// ─────────────────────────────────────────────
		// 5. Terminals & Fonts
		// ─────────────────────────────────────────────
		{
			ID:          "ghostty",
			Name:        "Ghostty",
			Description: "GPU-accelerated terminal emulator with native tabs and font features",
			Category:    CatTerminals,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "ghostty",
				distro.PMPacman: "ghostty",
			},
			FallbackRepo:    "ghostty-org/ghostty",
			BinaryName:      "ghostty",
			StowPackage:     "ghostty",
			DefaultSelected: true,
		},
		{
			ID:          "nerdfonts",
			Name:        "JetBrainsMono Nerd Font",
			Description: "Installs JetBrainsMono Nerd Font for powerline symbols and terminal icons",
			Category:    CatTerminals,
			Packages: map[distro.PackageManager]string{
				distro.PMPacman: "ttf-jetbrains-mono-nerd",
				distro.PMDnf:    "jetbrains-mono-fonts-all",
				distro.PMApt:    "fonts-jetbrains-mono",
				distro.PMBrew:   "font-jetbrains-mono-nerd-font",
			},
			CustomScript: "FONT_DIR=\"${HOME}/.local/share/fonts/JetBrainsMono\"; " +
				"mkdir -p \"${FONT_DIR}\"; " +
				"if [ ! -f \"${FONT_DIR}/JetBrainsMonoNerdFont-Regular.ttf\" ]; then " +
				"TMP_F=\"$(mktemp -d 2>/dev/null || mktemp -d -t 'font')\"; " +
				"curl -fsSL \"https://github.com/ryanoasis/nerd-fonts/releases/latest/download/JetBrainsMono.tar.xz\" -o \"${TMP_F}/font.tar.xz\" && " +
				"tar -xf \"${TMP_F}/font.tar.xz\" -C \"${FONT_DIR}\" 2>/dev/null || true; " +
				"rm -rf \"${TMP_F}\"; " +
				"if command -v fc-cache >/dev/null 2>&1; then fc-cache -f \"${FONT_DIR}\" >/dev/null 2>&1 || true; fi; " +
				"fi",
			DefaultSelected:  true,
			PostInstallNotes: "Installs into ~/.local/share/fonts and refreshes the font cache.",
		},

		// ─────────────────────────────────────────────
		// 6. Theming
		// ─────────────────────────────────────────────
		{
			ID:               "theme-catppuccin",
			Name:             "Catppuccin Preset",
			Description:      "Soothing pastel color scheme across terminal tools and editors",
			Category:         CatThemes,
			StowPackage:      "theme-catppuccin",
			DefaultSelected:  true,
			PostInstallNotes: "Configures Ghostty, Fish, Starship, Zellij, Micro, Bat, and Fastfetch with Catppuccin Mocha.",
		},
		{
			ID:               "theme-everforest",
			Name:             "Everforest Preset",
			Description:      "Warm, comfortable natural green color scheme across tools",
			Category:         CatThemes,
			StowPackage:      "theme-everforest",
			DefaultSelected:  false,
			PostInstallNotes: "Configures Ghostty, Fish, Starship, Zellij, Micro, Bat, and Fastfetch with Everforest Dark.",
		},
		{
			ID:               "theme-tokyonight",
			Name:             "Tokyo Night Preset",
			Description:      "Clean dark neon theme celebrating the lights of downtown Tokyo",
			Category:         CatThemes,
			StowPackage:      "theme-tokyonight",
			DefaultSelected:  false,
			PostInstallNotes: "Configures Ghostty, Fish, Starship, Zellij, Micro, Bat, and Fastfetch with Tokyo Night.",
		},

		// ─────────────────────────────────────────────
		// 7. Productivity & Cheatsheets
		// ─────────────────────────────────────────────
		{
			ID:          "fastfetch",
			Name:        "fastfetch",
			Description: "Fast, customizable neofetch alternative displaying system details",
			Category:    CatUtils,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "fastfetch",
				distro.PMPacman: "fastfetch",
				distro.PMDnf:    "fastfetch",
				distro.PMApt:    "fastfetch",
				distro.PMApk:    "fastfetch",
				distro.PMZypper: "fastfetch",
				distro.PMXbps:   "fastfetch",
				distro.PMBrew:   "fastfetch",
			},
			FallbackRepo:    "fastfetch-cli/fastfetch",
			BinaryName:      "fastfetch",
			StowPackage:     "fastfetch",
			DefaultSelected: true,
		},
		{
			ID:          "navi",
			Name:        "navi",
			Description: "Interactive cheatsheet searcher using fzf",
			Category:    CatUtils,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "navi",
				distro.PMPacman: "navi",
				distro.PMBrew:   "navi",
			},
			FallbackRepo:    "denisidoro/navi",
			BinaryName:      "navi",
			StowPackage:     "navi",
			DefaultSelected: false,
		},
		{
			ID:          "atuin",
			Name:        "Atuin",
			Description: "Searchable shell history backed by SQLite",
			Category:    CatUtils,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "atuin",
				distro.PMPacman: "atuin",
				distro.PMDnf:    "atuin",
				distro.PMApk:    "atuin",
				distro.PMBrew:   "atuin",
			},
			FallbackRepo:    "atuinsh/atuin",
			BinaryName:      "atuin",
			DefaultSelected: false,
		},
		{
			ID:          "xh",
			Name:        "xh",
			Description: "User-friendly HTTP client with concise syntax and colors",
			Category:    CatUtils,
			Packages: map[distro.PackageManager]string{
				distro.PMMoss:   "xh",
				distro.PMPacman: "xh",
				distro.PMApk:    "xh",
				distro.PMBrew:   "xh",
			},
			FallbackRepo:    "ducaale/xh",
			BinaryName:      "xh",
			DefaultSelected: false,
		},
	}
}
