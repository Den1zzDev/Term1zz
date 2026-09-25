# Contributing to Term1zz

Thank you for contributing to Term1zz. This guide explains how to add tools, update configurations, and test your changes across supported distributions.

## Architecture

Term1zz is structured into focused Go packages:

- `cmd/term1zz`: Main entry point and CLI flag parsing.
- `internal/distro`: Distribution detection, package manager mapping, and privilege escalation.
- `internal/registry`: Catalog of tools, categories, package mappings, and install scripts.
- `internal/configs`: Dotfile manager that creates symlinks and writes timestamped backups.
- `internal/runner`: Dependency resolution, theme override logic, execution planner, and progress reporting.
- `internal/tui`: Terminal user interface built with Bubbletea and Lipgloss.
- `stow/`: Configuration files organized by tool and theme preset.

## Adding a new tool

1. Open `internal/registry/items.go`.
2. Add an `Item` struct to `Items()` under the relevant category:

```go
{
    ID:          "mytool",
    Name:        "My Tool",
    Description: "Description of what the tool does",
    Category:    CatCoreutils,
    Packages: map[distro.PackageManager]string{
        distro.PMMoss:   "mytool",
        distro.PMPacman: "mytool",
        distro.PMDnf:    "mytool",
        distro.PMApt:    "mytool",
        distro.PMApk:    "mytool",
        distro.PMZypper: "mytool",
        distro.PMXbps:   "mytool",
        distro.PMBrew:   "mytool",
    },
    BinaryName:      "mytool",
    DefaultSelected: true,
}
```

3. If the tool includes dotfiles, set `StowPackage: "mytool"` and add the configuration files under `stow/mytool/`. Mirror the target structure relative to the user home directory (for example, `stow/mytool/.config/mytool/config.toml`).

## Adding or editing theme presets

Theme presets live under `stow/theme-<name>/`. Each theme provides configurations for all integrated tools:

- `stow/theme-<name>/.config/ghostty/config`
- `stow/theme-<name>/.config/starship.toml`
- `stow/theme-<name>/.config/zellij/config.kdl`
- `stow/theme-<name>/.config/micro/settings.json`
- `stow/theme-<name>/.config/bat/config`
- `stow/theme-<name>/.config/fastfetch/config.jsonc`
- `stow/theme-<name>/.config/fish/conf.d/00-term1zz-theme.fish`

When a user selects a theme preset, Term1zz applies the theme's configurations over the default tool configs.

## Development and testing

### Building locally

```sh
go build -o term1zz ./cmd/term1zz
```

### Static analysis

```sh
go vet ./...
```

### Testing across distributions

Never run destructive package install tests directly on your host machine. Use container testing scripts with Podman:

```sh
# Dry run inside Arch Linux
podman run --rm -v $(pwd):/repo:ro docker.io/library/archlinux:latest sh -c "
    pacman -Syu --noconfirm which sudo >/dev/null 2>&1
    useradd -m tester && echo 'tester ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers
    su - tester -c 'cp -r /repo ~/Term1zz && cd ~/Term1zz && ./setup.sh --batch --dry-run'
"

# Dry run inside Fedora
podman run --rm -v $(pwd):/repo:ro docker.io/library/fedora:latest sh -c "
    dnf install -y which sudo >/dev/null 2>&1
    useradd -m tester && echo 'tester ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers
    su - tester -c 'cp -r /repo ~/Term1zz && cd ~/Term1zz && ./setup.sh --batch --dry-run'
"

# Dry run inside Ubuntu
podman run --rm -v $(pwd):/repo:ro docker.io/library/ubuntu:latest sh -c "
    apt-get update -qq && apt-get install -y -qq sudo >/dev/null 2>&1
    useradd -m -s /bin/bash tester && echo 'tester ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers
    su - tester -c 'cp -r /repo ~/Term1zz && cd ~/Term1zz && ./setup.sh --batch --dry-run'
"

# Dry run inside Alpine
podman run --rm -v $(pwd):/repo:ro docker.io/library/alpine:latest sh -c "
    apk add --no-cache sudo bash >/dev/null 2>&1
    adduser -D tester && echo 'tester ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers
    su - tester -c 'cp -r /repo ~/Term1zz && cd ~/Term1zz && ./setup.sh --batch --dry-run'
"
```

## Pull requests

- Keep changes focused and self-contained.
- Verify that code compiles with `CGO_ENABLED=0 go build ./cmd/term1zz`.
- Ensure scripts run under standard `/bin/sh` without bashisms.
