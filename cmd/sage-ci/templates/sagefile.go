package main

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/workflows"
	"go.einride.tech/sage/sg"
)

// cfg defines the project-specific configuration for sage-ci.
// Customize this to match your project structure.
var cfg = workflows.Config{
	// GoModules lists the Go module paths relative to the repository root.
	// Example: []string{".", "tools"}
	GoModules: []string{},

	// PythonModules lists the Python module paths relative to the repository root.
	// Example: []string{".", "scripts"}
	PythonModules: []string{},

	// LuaModules lists the Lua module paths relative to the repository root.
	// Example: []string{".", "plugins"}
	LuaModules: []string{},

	// Skip lists workflow names to skip during sync.
	// Example: []string{"lint", "test"}
	Skip: []string{},

	// OutputDir specifies a custom output directory for workflows.
	// Defaults to ".github/workflows" if empty.
	OutputDir: "",
}

// skipTargets lists sage target names to skip in RunSynced and RunSyncedSerial.
// Key: Target name (e.g. "GoTest").
// Value: List of modules to skip. Use "*" or match the module name.
// Example: map[string][]string{"GoLint": {"tools"}}
var skipTargets = map[string][]string{}

func main() {
	sg.GenerateMakefiles(
		sg.Makefile{
			Path:          sg.FromGitRoot("Makefile"),
			DefaultTarget: All,
		},
	)
}

// All is the default target. Customize this to run the targets you need.
func All(ctx context.Context) error {
	sg.Deps(ctx, RunSyncedSerial)
	sg.Deps(ctx, RunSynced)
	return GitDiffCheck(ctx)
}
