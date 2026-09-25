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
				plan.CustomScripts = append(plan.CustomScripts, ScriptAction{
					Name:   sel.Item.Name,
					Script: sel.Item.CustomScript,
				})
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

	// Deduplicate ConfigPlan by TargetPath so the active theme overrides generic configs
	targetMap := make(map[string]configs.PlanItem)
	var orderedTargets []string
	for _, item := range allPlanItems {
		if _, exists := targetMap[item.TargetPath]; !exists {
			orderedTargets = append(orderedTargets, item.TargetPath)
		}
		targetMap[item.TargetPath] = item
	}
	for _, target := range orderedTargets {
		plan.ConfigPlan = append(plan.ConfigPlan, targetMap[target])
	}

	return plan, nil
}

// Execute runs the plan, streaming progress messages to the log channel.
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

	// Step 1: Package installation
	if len(p.PackagesToPM) > 0 {
		cmdParts := p.DistroInfo.InstallCommand(p.PackagesToPM)
		if len(cmdParts) == 0 {
			return fmt.Errorf("unable to determine package install command for %s", p.DistroInfo.PM)
		}

		logChan <- fmt.Sprintf("Running: %s", strings.Join(cmdParts, " "))
		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		cmd.Stderr = cmd.Stdout

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed starting package manager: %w", err)
		}

		var lastOutput []string
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			text := scanner.Text()
			logChan <- text
			lastOutput = append(lastOutput, text)
			if len(lastOutput) > 8 {
				lastOutput = lastOutput[1:]
			}
		}

		if err := cmd.Wait(); err != nil {
			errDetail := strings.Join(lastOutput, " | ")
			if errDetail != "" {
				return fmt.Errorf("package manager failed: %s (%w)", errDetail, err)
			}
			return fmt.Errorf("package manager exited with error: %w", err)
		}
		logChan <- "Package installation finished."
	}

	// Step 2: Dotfile deployment
	if len(p.ConfigPlan) > 0 {
		if p.ActiveThemeName != "" {
			logChan <- fmt.Sprintf("Deploying %s across tools (Fish, Starship, Zellij, Ghostty, Micro, Bat, Fastfetch)...", p.ActiveThemeName)
		}
		logChan <- fmt.Sprintf("Deploying %d configuration files...", len(p.ConfigPlan))
		if err := cfgMgr.Apply(p.ConfigPlan); err != nil {
			return fmt.Errorf("configuration deployment failed: %w", err)
		}
		logChan <- fmt.Sprintf("Configurations deployed. Backups saved to %s", cfgMgr.BackupDir)
	}

	// Step 3: Custom script hooks (e.g. font installer or fisher setup)
	if len(p.CustomScripts) > 0 {
		for _, sa := range p.CustomScripts {
			logChan <- fmt.Sprintf("Running custom setup for %s...", sa.Name)
			cmd := exec.Command("sh", "-c", sa.Script)
			out, err := cmd.CombinedOutput()
			if err != nil {
				logChan <- fmt.Sprintf("Warning: setup hook for %s reported: %v (%s)", sa.Name, err, strings.TrimSpace(string(out)))
			} else {
				logChan <- fmt.Sprintf("Setup completed for %s.", sa.Name)
			}
		}
	}

	logChan <- "Term1zz setup complete!"
	return nil
}
