package tui

import (
	"fmt"
	"strings"

	"github.com/Den1zzDev/Term1zz/internal/configs"
	"github.com/Den1zzDev/Term1zz/internal/distro"
	"github.com/Den1zzDev/Term1zz/internal/registry"
	"github.com/Den1zzDev/Term1zz/internal/runner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UIState int

const (
	StateSelecting UIState = iota
	StateReviewing
	StateRunning
	StateFinished
)

type logMsg string
type finishedMsg struct{ err error }

// Model represents the Bubbletea application state.
type Model struct {
	DistroInfo distro.Info
	CfgMgr     *configs.Manager
	Items      []registry.Item
	Categories []registry.Category

	ActiveCatIdx int
	Cursor       int
	Selected     map[string]bool
	Modes        map[string]runner.InstallMode

	State   UIState
	Plan    *runner.Plan
	Logs    []string
	Err     error
	DryRun  bool
	Width   int
	Height  int
	logChan chan string
}

// NewModel initializes the TUI model.
func NewModel(d distro.Info, cfgMgr *configs.Manager, dryRun bool, preferredTheme string) Model {
	items := registry.Items()
	cats := registry.AllCategories()

	selected := make(map[string]bool)
	modes := make(map[string]runner.InstallMode)

	for _, it := range items {
		selected[it.ID] = it.DefaultSelected
		if it.Category == registry.CatThemes && preferredTheme != "" {
			selected[it.ID] = (it.ID == preferredTheme)
		}
		if it.HasConfig() {
			modes[it.ID] = runner.ModeBoth
		} else {
			modes[it.ID] = runner.ModePackageOnly
		}
	}

	return Model{
		DistroInfo: d,
		CfgMgr:     cfgMgr,
		Items:      items,
		Categories: cats,
		Selected:   selected,
		Modes:      modes,
		State:      StateSelecting,
		DryRun:     dryRun,
		Width:      80,
		Height:     24,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// filteredItems returns the items belonging to the currently active category.
func (m Model) currentCategoryItems() []registry.Item {
	if m.ActiveCatIdx >= len(m.Categories) {
		return nil
	}
	currentCat := m.Categories[m.ActiveCatIdx]
	var res []registry.Item
	for _, it := range m.Items {
		if it.Category == currentCat {
			res = append(res, it)
		}
	}
	return res
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		switch m.State {
		case StateSelecting:
			return m.updateSelecting(msg)
		case StateReviewing:
			return m.updateReviewing(msg)
		case StateRunning:
			return m, nil // Don't interrupt running tasks with regular keys
		case StateFinished:
			if msg.String() == "q" || msg.String() == "enter" || msg.String() == "esc" {
				return m, tea.Quit
			}
		}

	case logMsg:
		m.Logs = append(m.Logs, string(msg))
		return m, waitForLog(m.logChan)

	case finishedMsg:
		m.State = StateFinished
		m.Err = msg.err
		return m, nil
	}

	return m, nil
}

func (m Model) updateSelecting(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	items := m.currentCategoryItems()

	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if m.Cursor < len(items)-1 {
			m.Cursor++
		}

	case "1", "2", "3", "4", "5", "6", "7":
		idx := int(msg.String()[0] - '1')
		if idx >= 0 && idx < len(m.Categories) {
			m.ActiveCatIdx = idx
			m.Cursor = 0
		}

	case "tab", "right", "l":
		m.ActiveCatIdx = (m.ActiveCatIdx + 1) % len(m.Categories)
		m.Cursor = 0

	case "shift+tab", "left", "h":
		m.ActiveCatIdx = (m.ActiveCatIdx - 1 + len(m.Categories)) % len(m.Categories)
		m.Cursor = 0

	case " ":
		if len(items) > 0 && m.Cursor < len(items) {
			it := items[m.Cursor]
			if it.Category == registry.CatThemes {
				// Mutually exclusive: choosing a theme activates it and unselects all other themes
				for _, other := range m.Items {
					if other.Category == registry.CatThemes {
						m.Selected[other.ID] = (other.ID == it.ID)
					}
				}
			} else {
				m.Selected[it.ID] = !m.Selected[it.ID]
			}
		}

	case "a":
		if m.ActiveCatIdx < len(m.Categories) && m.Categories[m.ActiveCatIdx] == registry.CatThemes {
			// In Theming tab, 'a' exclusively selects the currently highlighted theme
			if len(items) > 0 && m.Cursor < len(items) {
				it := items[m.Cursor]
				for _, other := range m.Items {
					if other.Category == registry.CatThemes {
						m.Selected[other.ID] = (other.ID == it.ID)
					}
				}
			}
			break
		}

		// Toggle all in category
		allSelected := true
		for _, it := range items {
			if !m.Selected[it.ID] {
				allSelected = false
				break
			}
		}
		for _, it := range items {
			m.Selected[it.ID] = !allSelected
		}

	case "m":
		// Cycle mode between Both, PackageOnly, ConfigOnly
		if len(items) > 0 && m.Cursor < len(items) {
			it := items[m.Cursor]
			if it.HasConfig() {
				curr := m.Modes[it.ID]
				m.Modes[it.ID] = (curr + 1) % 3
			}
		}

	case "d":
		// Toggle dry-run
		m.DryRun = !m.DryRun

	case "enter":
		// Build plan and move to review state
		var selections []runner.Selection
		for _, it := range m.Items {
			if m.Selected[it.ID] {
				selections = append(selections, runner.Selection{
					Item: it,
					Mode: m.Modes[it.ID],
				})
			}
		}

		plan, err := runner.BuildPlan(m.DistroInfo, m.CfgMgr, selections, m.DryRun)
		if err != nil {
			m.Err = err
			return m, nil
		}
		m.Plan = plan
		m.State = StateReviewing
	}

	return m, nil
}

func (m Model) updateReviewing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "b":
		m.State = StateSelecting
		return m, nil

	case "y", "enter":
		m.State = StateRunning
		m.logChan = make(chan string, 100)
		return m, tea.Batch(
			runPlanCmd(m.Plan, m.CfgMgr, m.logChan),
			waitForLog(m.logChan),
		)
	}
	return m, nil
}

func runPlanCmd(p *runner.Plan, cfgMgr *configs.Manager, logChan chan string) tea.Cmd {
	return func() tea.Msg {
		err := p.Execute(cfgMgr, logChan)
		return finishedMsg{err: err}
	}
}

func waitForLog(logChan <-chan string) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-logChan
		if !ok {
			return nil
		}
		return logMsg(msg)
	}
}

func modeLabel(m runner.InstallMode, hasConfig bool) string {
	if !hasConfig {
		return "[pkg]"
	}
	switch m {
	case runner.ModeBoth:
		return "[all]"
	case runner.ModePackageOnly:
		return "[pkg]"
	case runner.ModeConfigOnly:
		return "[cfg]"
	}
	return "[all]"
}

func (m Model) View() string {
	switch m.State {
	case StateSelecting:
		return m.viewSelecting()
	case StateReviewing:
		return m.viewReviewing()
	case StateRunning:
		return m.viewRunning()
	case StateFinished:
		return m.viewFinished()
	}
	return ""
}

func (m Model) viewSelecting() string {
	var b strings.Builder

	// Header Banner
	header := TitleStyle.Render("✦ Term1zz ✦") + " " + SubtleStyle.Render("Modular Terminal Toolbox")
	if m.DryRun {
		header += " " + BadgeStyle.Render("DRY RUN")
	}
	b.WriteString(header + "\n")
	b.WriteString(SubtleStyle.Render(m.DistroInfo.FormatSummary()) + "\n\n")

	// Category Tabs (clean two-line display with 1-7 numbers)
	var catTabsLine1 []string
	var catTabsLine2 []string
	for i, cat := range m.Categories {
		tabLabel := fmt.Sprintf("[%d] %s", i+1, cat)
		var rendered string
		if i == m.ActiveCatIdx {
			rendered = SelectedCategoryStyle.Render(tabLabel)
		} else {
			rendered = CategoryStyle.Render(tabLabel)
		}
		if i < 4 {
			catTabsLine1 = append(catTabsLine1, rendered)
		} else {
			catTabsLine2 = append(catTabsLine2, rendered)
		}
	}
	b.WriteString(strings.Join(catTabsLine1, " ") + "\n")
	b.WriteString(strings.Join(catTabsLine2, " ") + "\n\n")

	// Main Split: Item List on left, Details on right
	items := m.currentCategoryItems()

	var listLines []string
	for i, it := range items {
		check := UncheckedStyle.Render("[ ]")
		if m.Selected[it.ID] {
			check = CheckedStyle.Render("[✓]")
		}

		mLabel := SubtleStyle.Render(modeLabel(m.Modes[it.ID], it.HasConfig()))

		line := fmt.Sprintf("%s %s %s", check, mLabel, it.Name)
		if i == m.Cursor {
			line = SelectedItemStyle.Render("> " + line)
		} else {
			line = ItemStyle.Render("  " + line)
		}
		listLines = append(listLines, line)
	}

	// Detail Pane for currently focused item
	var detailContent string
	if len(items) > 0 && m.Cursor < len(items) {
		focused := items[m.Cursor]
		pkgName, found := focused.PackageFor(m.DistroInfo.PM)
		pkgStatus := pkgName
		if !found {
			if focused.CustomScript != "" {
				pkgStatus = "managed via custom setup script"
			} else if focused.StowPackage != "" {
				pkgStatus = "suite theme configuration preset"
			} else {
				pkgStatus = "no native package"
			}
		}

		detailContent = fmt.Sprintf("%s\n\n%s\n\n%s: %s\n%s: %t\n%s: %s",
			TitleStyle.Render(focused.Name),
			focused.Description,
			KeyStyle.Render("Package ("+string(m.DistroInfo.PM)+")"), pkgStatus,
			KeyStyle.Render("Dotfiles bundled"), focused.HasConfig(),
			KeyStyle.Render("Install Mode"), modeLabel(m.Modes[focused.ID], focused.HasConfig()),
		)

		if focused.PostInstallNotes != "" {
			detailContent += fmt.Sprintf("\n\n%s: %s", KeyStyle.Render("Suite Scope"), focused.PostInstallNotes)
		}
	}

	termWidth := m.Width
	if termWidth <= 0 {
		termWidth = 80
	}

	var split string
	if termWidth < 95 {
		// Responsive vertical layout for 80-column / smaller terminals
		contentWidth := termWidth - 4
		if contentWidth < 40 {
			contentWidth = 40
		}
		leftPane := ListBoxStyle.Width(contentWidth).Render(strings.Join(listLines, "\n"))
		rightPane := DetailBoxStyle.Width(contentWidth).Render(detailContent)
		split = lipgloss.JoinVertical(lipgloss.Left, leftPane, rightPane)
	} else {
		// Responsive side-by-side split for wide terminals
		leftWidth := (termWidth - 6) * 55 / 100
		if leftWidth < 44 {
			leftWidth = 44
		}
		rightWidth := termWidth - leftWidth - 6
		if rightWidth < 34 {
			rightWidth = 34
		}
		leftPane := ListBoxStyle.Width(leftWidth).Height(12).Render(strings.Join(listLines, "\n"))
		rightPane := DetailBoxStyle.Width(rightWidth).Height(12).Render(detailContent)
		split = lipgloss.JoinHorizontal(lipgloss.Top, leftPane, "  ", rightPane)
	}

	b.WriteString(split + "\n\n")

	// Keymap Footer
	footer := fmt.Sprintf("%s toggle  %s mode  %s all  %s category  %s dry-run  %s review  %s quit",
		KeyStyle.Render("[space]"),
		KeyStyle.Render("[m]"),
		KeyStyle.Render("[a]"),
		KeyStyle.Render("[tab]"),
		KeyStyle.Render("[d]"),
		KeyStyle.Render("[enter]"),
		KeyStyle.Render("[q]"),
	)
	b.WriteString(FooterStyle.Render(footer))

	return b.String()
}

func (m Model) viewReviewing() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("✦ Installation Plan Review ✦") + "\n\n")

	if m.Plan == nil {
		b.WriteString("No plan generated.\n")
		return b.String()
	}

	if m.DryRun {
		b.WriteString(BadgeStyle.Render("SIMULATION MODE ACTIVE (no changes will be made)") + "\n\n")
	}

	if m.Plan.ActiveThemeName != "" {
		b.WriteString(KeyStyle.Render("Active Suite Theme:\n"))
		b.WriteString(fmt.Sprintf("  ✦ %s\n", m.Plan.ActiveThemeName))
		b.WriteString(SubtleStyle.Render("    Configures Fish, Starship, Zellij, Ghostty, Micro, Bat, and Fastfetch\n\n"))
	}

	b.WriteString(KeyStyle.Render("Packages to install via " + string(m.DistroInfo.PM) + ":\n"))
	if len(m.Plan.PackagesToPM) == 0 {
		b.WriteString("  (None selected)\n")
	} else {
		b.WriteString("  " + strings.Join(m.Plan.PackagesToPM, ", ") + "\n")
	}
	b.WriteString("\n")

	b.WriteString(KeyStyle.Render("Configuration packages to deploy:\n"))
	if len(m.Plan.ConfigsToLink) == 0 {
		b.WriteString("  (None selected)\n")
	} else {
		b.WriteString("  " + strings.Join(m.Plan.ConfigsToLink, ", ") + "\n")
		b.WriteString(fmt.Sprintf("  Total files to link: %d\n", len(m.Plan.ConfigPlan)))
		b.WriteString(fmt.Sprintf("  Backups target: %s\n", m.CfgMgr.BackupDir))
	}
	b.WriteString("\n")

	if len(m.Plan.CustomScripts) > 0 {
		b.WriteString(KeyStyle.Render("Custom setup scripts to run:\n"))
		for _, sa := range m.Plan.CustomScripts {
			b.WriteString("  ▸ " + sa.Name + "\n")
		}
		b.WriteString("\n")
	}

	if len(m.Plan.Warnings) > 0 {
		b.WriteString(WarningStyle.Render("Warnings / Notes:\n"))
		for _, w := range m.Plan.Warnings {
			b.WriteString("  ! " + w + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(FooterStyle.Render(fmt.Sprintf("%s confirm and run  %s back  %s quit",
		KeyStyle.Render("[y / enter]"),
		KeyStyle.Render("[b / esc]"),
		KeyStyle.Render("[q]"),
	)))

	return b.String()
}

func (m Model) viewRunning() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("✦ Executing Term1zz Setup ✦") + "\n\n")

	recentLogs := m.Logs
	if len(recentLogs) > 15 {
		recentLogs = recentLogs[len(recentLogs)-15:]
	}

	b.WriteString(ListBoxStyle.Width(80).Height(16).Render(strings.Join(recentLogs, "\n")))
	b.WriteString("\n\n" + SubtleStyle.Render("Processing... please wait."))
	return b.String()
}

func (m Model) viewFinished() string {
	var b strings.Builder

	if m.Err != nil {
		b.WriteString(TitleStyle.Render("Setup Failed") + "\n\n")
		b.WriteString(WarningStyle.Render(fmt.Sprintf("Error: %v\n\n", m.Err)))
	} else {
		b.WriteString(TitleStyle.Render("✦ Setup Completed Successfully! ✦") + "\n\n")
		if m.DryRun {
			b.WriteString("Dry run completed cleanly. No files or packages were modified.\n\n")
		} else {
			b.WriteString("All selected tools and configurations have been set up.\n")
			b.WriteString(fmt.Sprintf("Any existing configurations were safely backed up to:\n%s\n\n", m.CfgMgr.BackupDir))
			b.WriteString("Restart your shell or terminal emulator to load new environments.\n\n")
		}
	}

	b.WriteString(FooterStyle.Render(KeyStyle.Render("[q / enter]") + " to exit"))
	return b.String()
}
