package main

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/targets"
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

// skip lists sage target names to skip.
// Key: Target name (e.g. "GoTest").
// Value: List of modules to skip. Use "*" to skip all modules.
// Example: targets.SkipTargets{"GoLint": {"tools"}}
var skip = targets.SkipTargets{}

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
	sg.Deps(ctx, Sync)
	sg.Deps(ctx, RunSerial)
	sg.Deps(ctx, RunParallel)
	return targets.GitDiffCheck(ctx)
}

// Sync regenerates GitHub Actions workflows.
func Sync(ctx context.Context) error {
	return targets.GenerateGHA(cfg)
}

// RunSerial runs all mutating targets serially.
func RunSerial(ctx context.Context) error {
	return targets.RunSerial(ctx, cfg, skip)
}

// RunParallel runs all non-mutating targets in parallel.
func RunParallel(ctx context.Context) error {
	return targets.RunParallel(ctx, cfg, skip)
}

// --- Individual targets (uncomment to expose in Makefile) ---

// func GoModTidy(ctx context.Context) error   { return targets.GoModTidy(ctx, cfg, skip) }
// func GoFormat(ctx context.Context) error    { return targets.GoFormat(ctx, cfg, skip) }
// func GoLint(ctx context.Context) error      { return targets.GoLint(ctx, cfg, skip) }
// func GoTest(ctx context.Context) error      { return targets.GoTest(ctx, cfg, skip) }
// func GoVulncheck(ctx context.Context) error { return targets.GoVulncheck(ctx, cfg, skip) }

// func PythonSync(ctx context.Context) error   { return targets.PythonSync(ctx, cfg, skip) }
// func PythonFormat(ctx context.Context) error { return targets.PythonFormat(ctx, cfg, skip) }
// func PythonLint(ctx context.Context) error   { return targets.PythonLint(ctx, cfg, skip) }
// func PythonMypy(ctx context.Context) error   { return targets.PythonMypy(ctx, cfg, skip) }
// func PythonTest(ctx context.Context) error   { return targets.PythonTest(ctx, cfg, skip) }

// func LuaFormat(ctx context.Context) error { return targets.LuaFormat(ctx, cfg, skip) }
