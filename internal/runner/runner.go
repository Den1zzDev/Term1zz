package runner

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Den1zzDev/Term1zz/internal/configs"
	"github.com/Den1zzDev/Term1zz/internal/distro"
	"github.com/Den1zzDev/Term1zz/internal/registry"
)

// InstallMode defines what to execute for an item.
type InstallMode int

const (
	ModeBoth InstallMode = iota
	ModePackageOnly
	ModeConfigOnly
)

// Selection represents a user's choice for a specific item.
type Selection struct {
	Item registry.Item
	Mode InstallMode
}

// ScriptAction represents an optional custom setup script to execute.
type ScriptAction struct {
	Name   string
	Script string
}

// Plan represents the full sequence of actions to be executed.
type Plan struct {
	DistroInfo      distro.Info
	ActiveThemeName string
	PackagesToPM    []string
	ConfigsToLink   []string
	ConfigPlan      []configs.PlanItem
	CustomScripts   []ScriptAction
	Warnings        []string
	DryRun          bool
}

// BuildPlan constructs a runnable plan from user selections.
func BuildPlan(d distro.Info, cfgMgr *configs.Manager, selections []Selection, dryRun bool) (*Plan, error) {
	plan := &Plan{
		DistroInfo: d,
		DryRun:     dryRun,
	}

	seenPkgs := make(map[string]bool)

	// Separate regular tools from active theme selection
	var themeSelection *Selection
	var regularSelections []Selection

	for i := range selections {
		sel := selections[i]
		if sel.Item.Category == registry.CatThemes {
			themeSelection = &sel
		} else {
			regularSelections = append(regularSelections, sel)
		}
	}

	if themeSelection != nil {
		plan.ActiveThemeName = themeSelection.Item.Name
	}

	// Process regular tools first, then the active theme last
	// so the theme's tailored configurations take precedence for all tools
	var orderedSelections []Selection
	orderedSelections = append(orderedSelections, regularSelections...)
	if themeSelection != nil {
		orderedSelections = append(orderedSelections, *themeSelection)
	}

	var allPlanItems []configs.PlanItem
	seenConfigs := make(map[string]bool)

	for _, sel := range orderedSelections {
		// Package resolution
		if sel.Mode == ModeBoth || sel.Mode == ModePackageOnly {
			pkg, found := sel.Item.PackageFor(d.PM)
			if found {
				if !seenPkgs[pkg] {
					seenPkgs[pkg] = true
					plan.PackagesToPM = append(plan.PackagesToPM, pkg)
				}
			} else if sel.Item.CustomScript == "" && sel.Item.StowPackage == "" {
				plan.Warnings = append(plan.Warnings, fmt.Sprintf("%s has no native package for %s", sel.Item.Name, d.PM))
			}

			if sel.Item.CustomScript != "" {
				// Only run fallback scripts if no native package was found
				if !(sel.Item.CustomScriptFallback && found) {
					plan.CustomScripts = append(plan.CustomScripts, ScriptAction{
						Name:   sel.Item.Name,
						Script: sel.Item.CustomScript,
					})
				}
			}
		}

		// Config resolution
		if sel.Mode == ModeBoth || sel.Mode == ModeConfigOnly {
			if sel.Item.HasConfig() {
				pkgName := sel.Item.StowPackage
				if !seenConfigs[pkgName] {
					seenConfigs[pkgName] = true
					plan.ConfigsToLink = append(plan.ConfigsToLink, pkgName)
					items, err := cfgMgr.PlanPackage(pkgName)
					if err != nil {
						plan.Warnings = append(plan.Warnings, fmt.Sprintf("failed planning dotfiles for %s: %v", pkgName, err))
					} else {
						allPlanItems = append(allPlanItems, items...)
					}
				}
			}
		}
	}

	// Deduplicate planned config items by target path
	// Later entries (such as the active theme) override earlier entries for the same target
	targetMap := make(map[string]configs.PlanItem)
	var orderedTargets []string
	for _, it := range allPlanItems {
		if _, exists := targetMap[it.TargetPath]; !exists {
			orderedTargets = append(orderedTargets, it.TargetPath)
		}
		targetMap[it.TargetPath] = it
	}

	var dedupedPlan []configs.PlanItem
	for _, t := range orderedTargets {
		dedupedPlan = append(dedupedPlan, targetMap[t])
	}

	plan.ConfigPlan = dedupedPlan
	return plan, nil
}

// Execute runs the actions defined in the plan, streaming progress lines to logChan.
func (p *Plan) Execute(cfgMgr *configs.Manager, logChan chan<- string) error {
	defer close(logChan)

	if p.DryRun {
		logChan <- "[DRY RUN] Simulating Term1zz setup..."
		if p.ActiveThemeName != "" {
			logChan <- fmt.Sprintf("[DRY RUN] Active Theme Preset: %s (configuring terminal suite)", p.ActiveThemeName)
		}
		if len(p.PackagesToPM) > 0 {
			cmd := p.DistroInfo.InstallCommand(p.PackagesToPM)
			logChan <- fmt.Sprintf("[DRY RUN] Would execute package install: %s", strings.Join(cmd, " "))
		}
		for _, item := range p.ConfigPlan {
			logChan <- fmt.Sprintf("[DRY RUN] Config action %s: %s -> %s", item.Action, item.SourcePath, item.TargetPath)
		}
		for _, sa := range p.CustomScripts {
			logChan <- fmt.Sprintf("[DRY RUN] Would execute custom hook for %s", sa.Name)
		}
		logChan <- "[DRY RUN] Complete. No changes made."
		return nil
	}

	if p.ActiveThemeName != "" {
		logChan <- fmt.Sprintf("Applying %s suite theme across terminal utilities...", p.ActiveThemeName)
	}

	// 1. Install packages via package manager
	if len(p.PackagesToPM) > 0 {
		cmdParts := p.DistroInfo.InstallCommand(p.PackagesToPM)
		if len(cmdParts) > 0 {
			logChan <- fmt.Sprintf("Running package manager: %s", strings.Join(cmdParts, " "))
			cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
			cmd.Env = append(cmd.Environ(), "DEBIAN_FRONTEND=noninteractive")

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				return fmt.Errorf("failed to open stdout pipe: %w", err)
			}
			cmd.Stderr = cmd.Stdout

			if err := cmd.Start(); err != nil {
				return fmt.Errorf("failed to start package install command: %w", err)
			}

			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				logChan <- scanner.Text()
			}

			if err := cmd.Wait(); err != nil {
				return fmt.Errorf("package installation command failed: %w", err)
			}
			logChan <- "Package installation completed."
		}
	}

	// 2. Deploy dotfiles
	if len(p.ConfigPlan) > 0 {
		logChan <- fmt.Sprintf("Deploying %d configuration files...", len(p.ConfigPlan))
		if err := cfgMgr.Apply(p.ConfigPlan); err != nil {
			return fmt.Errorf("dotfile deployment failed: %w", err)
		}
		logChan <- "Configurations applied successfully."
	}

	// 3. Run custom post-install scripts
	for _, sa := range p.CustomScripts {
		logChan <- fmt.Sprintf("Executing setup hook: %s...", sa.Name)
		cmd := exec.Command("sh", "-c", sa.Script)
		out, err := cmd.CombinedOutput()
		if len(out) > 0 {
			for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
				logChan <- "  " + line
			}
		}
		if err != nil {
			logChan <- fmt.Sprintf("Warning: setup hook for %s exited with error: %v", sa.Name, err)
		} else {
			logChan <- fmt.Sprintf("Completed setup hook for %s.", sa.Name)
		}
	}

	logChan <- "All Term1zz operations completed successfully!"
	return nil
}
