package github

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fredrikaverpil/sage-ci/config"
)

func TestSync(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sage-ci-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	cfg := config.Config{
		GoModules: []string{"."},
		OutputDir: tmpDir,
	}

	if err := Sync(cfg); err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	expectedFiles := []string{
		"sage-ci-go-ci.yml",
		"sage-ci-pr.yml",
		"sage-ci-release.yml",
		"sage-ci-sync.yml",
		"sage-ci-stale.yml",
	}

	for _, file := range expectedFiles {
		path := filepath.Join(tmpDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s does not exist", file)
		}
	}

	// Verify skipping
	tmpDir2, _ := os.MkdirTemp("", "sage-ci-test-skip-*")
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir2) })

	cfg.OutputDir = tmpDir2
	cfg.SkipWorkflows = []string{"sage-ci-stale", "sage-ci-release"}

	if err := Sync(cfg); err != nil {
		t.Fatalf("Sync with skip failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir2, "sage-ci-stale.yml")); !os.IsNotExist(err) {
		t.Error("sage-ci-stale.yml should have been skipped")
	}
	if _, err := os.Stat(filepath.Join(tmpDir2, "sage-ci-release.yml")); !os.IsNotExist(err) {
		t.Error("sage-ci-release.yml should have been skipped")
	}
	if _, err := os.Stat(filepath.Join(tmpDir2, "sage-ci-go-ci.yml")); os.IsNotExist(err) {
		t.Error("sage-ci-go-ci.yml should not have been skipped")
	}

	// Verify ecosystem-specific workflows are skipped when no modules configured
	tmpDir3, _ := os.MkdirTemp("", "sage-ci-test-ecosystem-*")
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir3) })

	cfgNoGo := config.Config{
		OutputDir: tmpDir3,
		// No GoModules, PythonModules, or LuaModules
	}

	if err := Sync(cfgNoGo); err != nil {
		t.Fatalf("Sync with no modules failed: %v", err)
	}

	// Ecosystem-specific workflows should not exist
	if _, err := os.Stat(filepath.Join(tmpDir3, "sage-ci-go-ci.yml")); !os.IsNotExist(err) {
		t.Error("sage-ci-go-ci.yml should not exist when GoModules is empty")
	}

	// Generic workflows should still exist
	if _, err := os.Stat(filepath.Join(tmpDir3, "sage-ci-pr.yml")); os.IsNotExist(err) {
		t.Error("sage-ci-pr.yml should exist even without modules")
	}
}
