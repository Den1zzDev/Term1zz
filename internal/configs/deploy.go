package configs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ActionType identifies what the configuration manager will do with a file.
type ActionType string

const (
	ActionLink   ActionType = "LINK"
	ActionSkip   ActionType = "SKIP"
	ActionBackup ActionType = "BACKUP"
)

// PlanItem describes a planned file operation.
type PlanItem struct {
	SourcePath string
	TargetPath string
	Action     ActionType
	BackupPath string
	Detail     string
}

// Manager handles deploying curated dotfiles.
type Manager struct {
	StowDir   string
	HomeDir   string
	BackupDir string
	DryRun    bool
}

// New creates a configuration manager.
func New(stowDir, homeDir string, dryRun bool) *Manager {
	timestamp := time.Now().Format("20060102-150405")
	backupDir := filepath.Join(homeDir, ".local", "state", "term1zz", "backups", timestamp)
	return &Manager{
		StowDir:   stowDir,
		HomeDir:   homeDir,
		BackupDir: backupDir,
		DryRun:    dryRun,
	}
}

// PlanPackage inspects a stow package and returns planned actions without touching the disk.
func (m *Manager) PlanPackage(pkgName string) ([]PlanItem, error) {
	pkgRoot := filepath.Join(m.StowDir, pkgName)
	if _, err := os.Stat(pkgRoot); err != nil {
		return nil, fmt.Errorf("package directory %s not found: %w", pkgRoot, err)
	}

	var plan []PlanItem

	err := filepath.Walk(pkgRoot, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(pkgRoot, srcPath)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(m.HomeDir, relPath)
		item := PlanItem{
			SourcePath: srcPath,
			TargetPath: targetPath,
		}

		targetInfo, err := os.Lstat(targetPath)
		if os.IsNotExist(err) {
			item.Action = ActionLink
			item.Detail = "New link"
			plan = append(plan, item)
			return nil
		}
		if err != nil {
			return err
		}

		// If target is already a symlink, check destination
		if targetInfo.Mode()&os.ModeSymlink != 0 {
			dest, err := os.Readlink(targetPath)
			if err == nil && (dest == srcPath || filepath.Clean(dest) == filepath.Clean(srcPath)) {
				item.Action = ActionSkip
				item.Detail = "Already linked correctly"
				plan = append(plan, item)
				return nil
			}
		}

		// Existing file or directory or different symlink requires backup
		item.Action = ActionBackup
		item.BackupPath = filepath.Join(m.BackupDir, pkgName, relPath)
		item.Detail = "Will backup existing file and link"
		plan = append(plan, item)
		return nil
	})

	return plan, err
}

// Apply executes a planned list of configuration operations.
func (m *Manager) Apply(plan []PlanItem) error {
	if m.DryRun {
		return nil
	}

	for _, item := range plan {
		switch item.Action {
		case ActionSkip:
			continue

		case ActionBackup:
			if err := m.backupFile(item.TargetPath, item.BackupPath); err != nil {
				return fmt.Errorf("failed backing up %s: %w", item.TargetPath, err)
			}
			if err := os.Remove(item.TargetPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed removing old file %s: %w", item.TargetPath, err)
			}
			fallthrough

		case ActionLink:
			if err := os.MkdirAll(filepath.Dir(item.TargetPath), 0755); err != nil {
				return fmt.Errorf("failed creating target dir: %w", err)
			}
			if err := os.Symlink(item.SourcePath, item.TargetPath); err != nil {
				return fmt.Errorf("failed creating symlink %s -> %s: %w", item.TargetPath, item.SourcePath, err)
			}
		}
	}
	return nil
}

func (m *Manager) backupFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	info, err := os.Lstat(src)
	if err != nil {
		return err
	}

	// If it's a symlink, recreate the symlink in the backup folder
	if info.Mode()&os.ModeSymlink != 0 {
		linkTarget, err := os.Readlink(src)
		if err != nil {
			return err
		}
		return os.Symlink(linkTarget, dst)
	}

	// Copy regular file
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
