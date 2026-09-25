package distro

import "fmt"

// InstallCommand builds the command slice to install one or more packages non-interactively.
func (i Info) InstallCommand(packages []string) []string {
	if len(packages) == 0 {
		return nil
	}

	var cmd []string
	if !i.IsRoot && i.Escalator != "" {
		if i.Escalator == "run0" {
			// run0 --background="" ensures no colored background block in modern systemd
			cmd = append(cmd, "run0", "--background=")
		} else {
			cmd = append(cmd, i.Escalator)
		}
	}

	switch i.PM {
	case PMMoss:
		cmd = append(cmd, "moss", "install", "-y")
		cmd = append(cmd, packages...)
	case PMPacman:
		cmd = append(cmd, "pacman", "-S", "--needed", "--noconfirm")
		cmd = append(cmd, packages...)
	case PMDnf:
		cmd = append(cmd, "dnf", "install", "-y")
		cmd = append(cmd, packages...)
	case PMApt:
		cmd = append(cmd, "apt-get", "install", "-y")
		cmd = append(cmd, packages...)
	case PMApk:
		cmd = append(cmd, "apk", "add")
		cmd = append(cmd, packages...)
	case PMZypper:
		cmd = append(cmd, "zypper", "--non-interactive", "install")
		cmd = append(cmd, packages...)
	case PMXbps:
		cmd = append(cmd, "xbps-install", "-y")
		cmd = append(cmd, packages...)
	case PMBrew:
		cmd = append(cmd, "brew", "install")
		cmd = append(cmd, packages...)
	default:
		return nil
	}

	return cmd
}

// FormatSummary returns a human-readable summary of the detected platform.
func (i Info) FormatSummary() string {
	name := i.DistroName
	if name == "" {
		name = i.DistroID
	}
	if name == "" {
		name = i.OS
	}
	return fmt.Sprintf("%s (%s) | Package manager: %s | Elevation: %s", name, i.Arch, i.PM, i.Escalator)
}
