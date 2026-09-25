package distro

import (
	"strings"
	"testing"
)

func TestDetect(t *testing.T) {
	info := Detect()
	if info.DistroID == "" {
		t.Error("expected non-empty DistroID")
	}
	if info.PM == "" {
		t.Error("expected non-empty PackageManager")
	}
	if info.HomeDir == "" {
		t.Error("expected non-empty HomeDir")
	}
	summary := info.FormatSummary()
	if !strings.Contains(summary, "Package manager:") {
		t.Errorf("summary missing 'Package manager:': %s", summary)
	}
}

func TestInstallCommand(t *testing.T) {
	pkgs := []string{"git", "curl"}

	tests := []struct {
		name       string
		info       Info
		wantCmd    string
		wantNoSudo bool
	}{
		{
			name:    "pacman non-root",
			info:    Info{PM: PMPacman, Escalator: "sudo", IsRoot: false},
			wantCmd: "sudo pacman -S --needed --noconfirm git curl",
		},
		{
			name:    "pacman root",
			info:    Info{PM: PMPacman, Escalator: "sudo", IsRoot: true},
			wantCmd: "pacman -S --needed --noconfirm git curl",
		},
		{
			name:    "moss non-root",
			info:    Info{PM: PMMoss, Escalator: "run0", IsRoot: false},
			wantCmd: "run0 --background= moss install -y git curl",
		},
		{
			name:       "brew always without escalator",
			info:       Info{PM: PMBrew, Escalator: "sudo", IsRoot: false},
			wantCmd:    "brew install git curl",
			wantNoSudo: true,
		},
		{
			name:    "brew with cask args",
			info:    Info{PM: PMBrew, Escalator: "", IsRoot: false},
			wantCmd: "brew install --cask ghostty",
		},
		{
			name:    "apt non-root",
			info:    Info{PM: PMApt, Escalator: "sudo", IsRoot: false},
			wantCmd: "sudo apt-get install -y git curl",
		},
		{
			name:    "dnf non-root",
			info:    Info{PM: PMDnf, Escalator: "sudo", IsRoot: false},
			wantCmd: "sudo dnf install -y git curl",
		},
		{
			name:    "apk non-root",
			info:    Info{PM: PMApk, Escalator: "doas", IsRoot: false},
			wantCmd: "doas apk add git curl",
		},
		{
			name:    "zypper non-root",
			info:    Info{PM: PMZypper, Escalator: "sudo", IsRoot: false},
			wantCmd: "sudo zypper --non-interactive install git curl",
		},
		{
			name:    "xbps non-root",
			info:    Info{PM: PMXbps, Escalator: "sudo", IsRoot: false},
			wantCmd: "sudo xbps-install -y git curl",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testPkgs []string
			if strings.Contains(tt.name, "cask") {
				testPkgs = []string{"--cask ghostty"}
			} else {
				testPkgs = pkgs
			}

			cmd := tt.info.InstallCommand(testPkgs)
			joined := strings.Join(cmd, " ")
			if joined != tt.wantCmd {
				t.Errorf("got %q, want %q", joined, tt.wantCmd)
			}
			if tt.wantNoSudo && strings.HasPrefix(joined, "sudo") {
				t.Errorf("brew command must not use sudo: %s", joined)
			}
		})
	}
}
