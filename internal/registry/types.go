package registry

import "github.com/Den1zzDev/Term1zz/internal/distro"

// Category represents a grouping of terminal components.
type Category string

const (
	CatShells      Category = "Shells & Prompts"
	CatMultiplex   Category = "Multiplexers"
	CatCoreutils   Category = "CLI Core Replacements"
	CatEditorsGit  Category = "Editors & Git"
	CatTerminals   Category = "Terminals & Fonts"
	CatThemes      Category = "Theming"
	CatUtils       Category = "Productivity & Cheatsheets"
)

// Item represents a tool, shell, font, or configuration package.
type Item struct {
	ID                   string
	Name                 string
	Description          string
	Category             Category
	Packages             map[distro.PackageManager]string
	FallbackRepo         string // GitHub owner/repo for direct releases
	BinaryName           string // Target executable name
	StowPackage          string // Directory under stow/ if curated dotfiles exist
	CustomScript         string // Optional shell command to run during installation
	CustomScriptFallback bool   // If true, CustomScript only runs when PackageFor finds no native package
	DefaultSelected      bool
	PostInstallNotes     string
}

// PackageFor returns the package name for the detected package manager.
func (it Item) PackageFor(pm distro.PackageManager) (string, bool) {
	if it.Packages == nil {
		return "", false
	}
	pkg, ok := it.Packages[pm]
	return pkg, ok && pkg != ""
}

// HasConfig returns true if Term1zz provides curated dotfiles for this item.
func (it Item) HasConfig() bool {
	return it.StowPackage != ""
}
