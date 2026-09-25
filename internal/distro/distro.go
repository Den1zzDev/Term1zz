package distro

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type PackageManager string

const (
	PMMoss    PackageManager = "moss"
	PMPacman  PackageManager = "pacman"
	PMDnf     PackageManager = "dnf"
	PMApt     PackageManager = "apt"
	PMApk     PackageManager = "apk"
	PMZypper  PackageManager = "zypper"
	PMXbps    PackageManager = "xbps"
	PMBrew    PackageManager = "brew"
	PMUnknown PackageManager = "unknown"
)

// Info holds detected distribution and environment data.
type Info struct {
	OS         string
	Arch       string
	DistroID   string
	DistroName string
	PM         PackageManager
	Escalator  string // "run0", "sudo", "doas", or empty if root / brew
	IsRoot     bool
	HomeDir    string
	ConfigDir  string
}

// Detect inspects the host environment.
func Detect() Info {
	info := Info{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		PM:        PMUnknown,
		IsRoot:    os.Geteuid() == 0,
		HomeDir:   os.Getenv("HOME"),
		ConfigDir: os.Getenv("XDG_CONFIG_HOME"),
	}

	if info.ConfigDir == "" && info.HomeDir != "" {
		info.ConfigDir = info.HomeDir + "/.config"
	}

	if info.OS == "darwin" {
		info.DistroID = "macos"
		info.DistroName = "macOS"
		ensureDarwinBrewPath()
		info.PM = detectPackageManager("macos")
		info.Escalator = "" // Homebrew prohibits running with sudo
	} else {
		parseOSRelease(&info)
		info.PM = detectPackageManager(info.DistroID)
		info.Escalator = detectEscalator(info.DistroID, info.IsRoot)
	}

	return info
}

func ensureDarwinBrewPath() {
	path := os.Getenv("PATH")
	for _, bp := range []string{"/opt/homebrew/bin", "/usr/local/bin"} {
		if _, err := os.Stat(bp + "/brew"); err == nil {
			if !strings.Contains(path, bp) {
				path = bp + ":" + path
				_ = os.Setenv("PATH", path)
			}
		}
	}
}

func parseOSRelease(info *Info) {
	paths := []string{"/etc/os-release", "/usr/lib/os-release"}
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := parts[0]
			val := strings.Trim(parts[1], `"'`)

			switch key {
			case "ID":
				info.DistroID = strings.ToLower(val)
			case "PRETTY_NAME":
				info.DistroName = val
			case "NAME":
				if info.DistroName == "" {
					info.DistroName = val
				}
			}
		}
		return
	}
}

func detectPackageManager(distroID string) PackageManager {
	// First check distro ID hints
	switch distroID {
	case "macos":
		if pathExists("brew") || pathExists("/opt/homebrew/bin/brew") || pathExists("/usr/local/bin/brew") {
			return PMBrew
		}
	case "aerynos", "serpent":
		if pathExists("moss") {
			return PMMoss
		}
	case "arch", "endeavouros", "manjaro", "artix", "cachyos":
		if pathExists("pacman") {
			return PMPacman
		}
	case "fedora", "rhel", "centos", "nobara":
		if pathExists("dnf") {
			return PMDnf
		}
	case "debian", "ubuntu", "pop", "linuxmint", "kali":
		if pathExists("apt-get") || pathExists("apt") {
			return PMApt
		}
	case "alpine":
		if pathExists("apk") {
			return PMApk
		}
	case "opensuse", "opensuse-tumbleweed", "opensuse-leap":
		if pathExists("zypper") {
			return PMZypper
		}
	case "void":
		if pathExists("xbps-install") {
			return PMXbps
		}
	}

	// Fallback to binary detection
	order := []struct {
		pm  PackageManager
		bin string
	}{
		{PMMoss, "moss"},
		{PMPacman, "pacman"},
		{PMDnf, "dnf"},
		{PMApt, "apt-get"},
		{PMApk, "apk"},
		{PMZypper, "zypper"},
		{PMXbps, "xbps-install"},
		{PMBrew, "brew"},
	}

	for _, check := range order {
		if pathExists(check.bin) {
			return check.pm
		}
	}

	return PMUnknown
}

func detectEscalator(distroID string, isRoot bool) string {
	if isRoot {
		return ""
	}
	// AerynOS uses run0 natively
	if (distroID == "aerynos" || distroID == "serpent") && pathExists("run0") {
		return "run0"
	}
	// For standard distros, containers, and general Linux, sudo is standard
	if pathExists("sudo") {
		return "sudo"
	}
	// Fall back to run0 only if systemd init is active
	if pathExists("run0") {
		if _, err := os.Stat("/run/systemd/system"); err == nil {
			return "run0"
		}
	}
	if pathExists("doas") {
		return "doas"
	}
	if pathExists("run0") {
		return "run0"
	}
	return ""
}

func pathExists(bin string) bool {
	_, err := exec.LookPath(bin)
	return err == nil
}
