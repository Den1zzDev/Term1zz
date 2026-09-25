package configs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPlanAndApply(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "term1zz-test-*")
	if err != nil {
		t.Fatalf("failed creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	stowDir := filepath.Join(tmpDir, "stow")
	homeDir := filepath.Join(tmpDir, "home")

	// Create mock stow package: mypkg/.config/app/config.toml
	pkgDir := filepath.Join(stowDir, "mypkg", ".config", "app")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatalf("failed creating mock stow package dir: %v", err)
	}
	sourceFile := filepath.Join(pkgDir, "config.toml")
	if err := os.WriteFile(sourceFile, []byte("key = 'value'\n"), 0644); err != nil {
		t.Fatalf("failed writing source file: %v", err)
	}

	mgr := New(stowDir, homeDir, false)

	// Plan the package
	plan, err := mgr.PlanPackage("mypkg")
	if err != nil {
		t.Fatalf("PlanPackage failed: %v", err)
	}
	if len(plan) != 1 {
		t.Fatalf("expected 1 plan item, got %d", len(plan))
	}

	item := plan[0]
	if item.Action != ActionLink {
		t.Errorf("expected ActionLink, got %s", item.Action)
	}
	expectedTarget := filepath.Join(homeDir, ".config", "app", "config.toml")
	if item.TargetPath != expectedTarget {
		t.Errorf("got target %s, want %s", item.TargetPath, expectedTarget)
	}

	// Apply plan
	if err := mgr.Apply(plan); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	// Verify symlink was created
	fi, err := os.Lstat(expectedTarget)
	if err != nil {
		t.Fatalf("target was not created: %v", err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Errorf("target is not a symlink")
	}

	// Verify idempotency: re-planning should result in ActionSkip
	plan2, err := mgr.PlanPackage("mypkg")
	if err != nil {
		t.Fatalf("second PlanPackage failed: %v", err)
	}
	if len(plan2) != 1 || plan2[0].Action != ActionSkip {
		t.Errorf("expected ActionSkip on second run, got %v", plan2[0].Action)
	}
}

func TestApplyBackupExisting(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "term1zz-backup-test-*")
	if err != nil {
		t.Fatalf("failed creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	stowDir := filepath.Join(tmpDir, "stow")
	homeDir := filepath.Join(tmpDir, "home")

	// Create mock stow package
	pkgDir := filepath.Join(stowDir, "pkg1", ".config")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(pkgDir, "settings.json"), []byte(`{"new":true}`), 0644)

	// Pre-create regular file in destination
	targetDir := filepath.Join(homeDir, ".config")
	_ = os.MkdirAll(targetDir, 0755)
	targetFile := filepath.Join(targetDir, "settings.json")
	_ = os.WriteFile(targetFile, []byte(`{"existing":true}`), 0644)

	mgr := New(stowDir, homeDir, false)
	plan, err := mgr.PlanPackage("pkg1")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan) != 1 || plan[0].Action != ActionBackup {
		t.Fatalf("expected ActionBackup, got %v", plan[0].Action)
	}

	if err := mgr.Apply(plan); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	// Check that backup directory contains the original file
	entries, err := os.ReadDir(mgr.BackupDir)
	if err != nil || len(entries) == 0 {
		t.Errorf("expected backup file in %s", mgr.BackupDir)
	}
}
