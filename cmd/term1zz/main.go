package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Den1zzDev/Term1zz/internal/configs"
	"github.com/Den1zzDev/Term1zz/internal/distro"
	"github.com/Den1zzDev/Term1zz/internal/registry"
	"github.com/Den1zzDev/Term1zz/internal/runner"
	"github.com/Den1zzDev/Term1zz/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

const version = "0.2.0"

func main() {
	dryRun := flag.Bool("dry-run", false, "Preview actions without installing packages or modifying dotfiles")
	flag.BoolVar(dryRun, "d", false, "Alias for --dry-run")

	batch := flag.Bool("batch", false, "Run in batch mode without interactive TUI")
	flag.BoolVar(batch, "b", false, "Alias for --batch")

	stowDir := flag.String("stow-dir", "", "Path to the stow configurations directory")
	flag.StringVar(stowDir, "s", "", "Alias for --stow-dir")

	themeFlag := flag.String("theme", "", "Suite theme to apply: catppuccin, everforest, or tokyonight")
	flag.StringVar(themeFlag, "t", "", "Alias for --theme")

	showVersion := flag.Bool("version", false, "Display Term1zz version")
	flag.BoolVar(showVersion, "v", false, "Alias for --version")

	flag.Parse()

	if *showVersion {
		fmt.Printf("Term1zz v%s\n", version)
		os.Exit(0)
	}

	distroInfo := distro.Detect()

	// Locate stow directory
	resolvedStowDir := findStowDir(*stowDir, distroInfo.HomeDir)
	cfgMgr := configs.New(resolvedStowDir, distroInfo.HomeDir, *dryRun)

	normalizedTheme := strings.ToLower(strings.TrimSpace(*themeFlag))
	if normalizedTheme != "" && !strings.HasPrefix(normalizedTheme, "theme-") {
		normalizedTheme = "theme-" + normalizedTheme
	}

	if *batch {
		runBatch(distroInfo, cfgMgr, *dryRun, normalizedTheme)
		return
	}

	m := tui.NewModel(distroInfo, cfgMgr, *dryRun, normalizedTheme)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running Term1zz TUI: %v\n", err)
		os.Exit(1)
	}
}

func findStowDir(custom, home string) string {
	if custom != "" {
		return custom
	}
	candidates := []string{
		"./stow",
		filepath.Join(home, ".local", "share", "Term1zz", "stow"),
		filepath.Join(home, ".local", "share", "term1zz", "stow"),
		"/usr/share/term1zz/stow",
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			abs, err := filepath.Abs(c)
			if err == nil {
				return abs
			}
			return c
		}
	}
	return "./stow"
}

func runBatch(d distro.Info, cfgMgr *configs.Manager, dryRun bool, activeTheme string) {
	fmt.Printf("Term1zz v%s (Batch Mode)\n", version)
	fmt.Println(d.FormatSummary())

	var selections []runner.Selection
	for _, it := range registry.Items() {
		isSelected := it.DefaultSelected
		if it.Category == registry.CatThemes && activeTheme != "" {
			isSelected = (it.ID == activeTheme)
		}

		if isSelected {
			mode := runner.ModeBoth
			if !it.HasConfig() {
				mode = runner.ModePackageOnly
			}
			selections = append(selections, runner.Selection{
				Item: it,
				Mode: mode,
			})
		}
	}

	plan, err := runner.BuildPlan(d, cfgMgr, selections, dryRun)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed building plan: %v\n", err)
		os.Exit(1)
	}

	logChan := make(chan string, 100)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for line := range logChan {
			fmt.Println(line)
		}
	}()

	if err := plan.Execute(cfgMgr, logChan); err != nil {
		wg.Wait()
		fmt.Fprintf(os.Stderr, "Batch execution failed: %v\n", err)
		os.Exit(1)
	}
	wg.Wait()
}
