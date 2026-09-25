package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Catppuccin Mocha / Frappe palette
	ColorMauve   = lipgloss.Color("#cba6f7")
	ColorPink    = lipgloss.Color("#f5c2e7")
	ColorGreen   = lipgloss.Color("#a6e3a1")
	ColorPeach   = lipgloss.Color("#fab387")
	ColorBlue    = lipgloss.Color("#89b4fa")
	ColorSubtext = lipgloss.Color("#a6adc8")
	ColorOverlay = lipgloss.Color("#6c7086")
	ColorSurface = lipgloss.Color("#313244")
	ColorBase    = lipgloss.Color("#1e1e2e")
	ColorText    = lipgloss.Color("#cdd6f4")
	ColorWarning = lipgloss.Color("#f38ba8")

	// Component styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorMauve).
			Padding(0, 1)

	BadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBase).
			Background(ColorMauve).
			Padding(0, 1)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(ColorOverlay)

	SelectedCategoryStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorMauve).
				Background(ColorSurface).
				Padding(0, 1)

	CategoryStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			Padding(0, 1)

	SelectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorText).
				Background(ColorSurface)

	ItemStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	CheckedStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	UncheckedStyle = lipgloss.NewStyle().
			Foreground(ColorOverlay)

	DetailBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorMauve).
			Padding(1, 2)

	ListBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorOverlay).
			Padding(0, 1)

	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			Padding(0, 1)

	KeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPeach)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning)
)
