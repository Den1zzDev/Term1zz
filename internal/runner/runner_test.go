package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Den1zzDev/Term1zz/internal/configs"
	"github.com/Den1zzDev/Term1zz/internal/distro"
	"github.com/Den1zzDev/Term1zz/internal/registry"
)

func TestBuildPlanThemePrecedence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "runner-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	stowDir := filepath.Join(tmpDir, "stow")
	homeDir := filepath.Join(tmpDir, "home")

	// Create tool stow package
	toolPkg := filepath.Join(stowDir, "ghostty", ".config", "ghostty")
	_ = os.MkdirAll(toolPkg, 0755)
	_ = os.WriteFile(filepath.Join(toolPkg, "config"), []byte("tool-config"), 0644)

	// Create theme stow package with same target
	themePkg := filepath.Join(stowDir, "theme-catppuccin", ".config", "ghostty")
	_ = os.MkdirAll(themePkg, 0755)
	_ = os.WriteFile(filepath.Join(themePkg, "config"), []byte("theme-config"), 0644)

	cfgMgr := configs.New(stowDir, homeDir, true)
	d := distro.Info{PM: distro.PMPacman, DistroID: "arch"}

	selections := []Selection{
		{
			Item: registry.Item{
				ID:          "ghostty",
				Name:        "Ghostty",
				StowPackage: "ghostty",
			},
			Mode: ModeConfigOnly,
		},
		{
			Item: registry.Item{
				ID:          "theme-catppuccin",
				Name:        "Catppuccin Preset",
				Category:    registry.CatThemes,
				StowPackage: "theme-catppuccin",
			},
			Mode: ModeConfigOnly,
		},
	}

	plan, err := BuildPlan(d, cfgMgr, selections, true)
	if err != nil {
		t.Fatalf("BuildPlan failed: %v", err)
	}

	if plan.ActiveThemeName != "Catppuccin Preset" {
		t.Errorf("got active theme %q, want %q", plan.ActiveThemeName, "Catppuccin Preset")
	}

	// Verify that deduplication kept the theme's source path
	if len(plan.ConfigPlan) != 1 {
		t.Fatalf("expected 1 config item after deduplication, got %d", len(plan.ConfigPlan))
	}
	expectedSource := filepath.Join(themePkg, "config")
	if plan.ConfigPlan[0].SourcePath != expectedSource {
		t.Errorf("theme did not take precedence; got source %s, want %s", plan.ConfigPlan[0].SourcePath, expectedSource)
	}
}

func TestBuildPlanCustomScriptFallback(t *testing.T) {
	dWithPacman := distro.Info{PM: distro.PMPacman, DistroID: "arch"}
	dWithMoss := distro.Info{PM: distro.PMMoss, DistroID: "aerynos"}
	cfgMgr := configs.New("/stow", "/home/test", true)

	item := registry.Item{
		ID:   "nerdfonts",
		Name: "JetBrainsMono Nerd Font",
		Packages: map[distro.PackageManager]string{
			distro.PMPacman: "ttf-jetbrains-mono-nerd",
		},
		CustomScript:         "curl -fsSL font.tar.xz",
		CustomScriptFallback: true,
	}

	// 1. On Arch where pacman package is found: CustomScript must NOT be scheduled
	planArch, err := BuildPlan(dWithPacman, cfgMgr, []Selection{{Item: item, Mode: ModeBoth}}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(planArch.PackagesToPM) != 1 || planArch.PackagesToPM[0] != "ttf-jetbrains-mono-nerd" {
		t.Errorf("expected pacman package, got %v", planArch.PackagesToPM)
	}
	if len(planArch.CustomScripts) != 0 {
		t.Errorf("expected 0 custom scripts when package exists, got %d", len(planArch.CustomScripts))
	}

	// 2. On Moss where package is missing: CustomScript MUST be scheduled
	planMoss, err := BuildPlan(dWithMoss, cfgMgr, []Selection{{Item: item, Mode: ModeBoth}}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(planMoss.PackagesToPM) != 0 {
		t.Errorf("expected no package for moss, got %v", planMoss.PackagesToPM)
	}
	if len(planMoss.CustomScripts) != 1 {
		t.Fatalf("expected fallback custom script to be scheduled, got %d", len(planMoss.CustomScripts))
	}
}

func TestDryRunExecution(t *testing.T) {
	d := distro.Info{PM: distro.PMPacman, DistroID: "arch"}
	cfgMgr := configs.New("/stow", "/home/test", true)

	plan := &Plan{
		DistroInfo:   d,
		DryRun:       true,
		PackagesToPM: []string{"git"},
		CustomScripts: []ScriptAction{
			{Name: "test-hook", Script: "exit 1"},
		},
	}

	logChan := make(chan string, 100)
	err := plan.Execute(cfgMgr, logChan)
	if err != nil {
		t.Fatalf("Dry run execution failed: %v", err)
	}

	var logs []string
	for l := range logChan {
		logs = append(logs, l)
	}

	if len(logs) == 0 {
		t.Error("expected dry run log lines")
	}
}
