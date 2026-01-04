package main

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/config"
	"github.com/fredrikaverpil/sage-ci/targets"
	"go.einride.tech/sage/sg"
)

// cfg defines the project-specific configuration for sage-ci.
// Customize this to match your project structure.
var cfg = config.Config{
	// GoModules lists the Go module paths relative to the repository root.
	// Example: []string{".", "tools"}
	GoModules: []string{},

	// PythonModules lists the Python module paths relative to the repository root.
	// Example: []string{".", "scripts"}
	PythonModules: []string{},

	// LuaModules lists the Lua module paths relative to the repository root.
	// Example: []string{".", "plugins"}
	LuaModules: []string{},

	// Platform specifies which CI platform to generate workflows for.
	// Options: "github", "gitlab", "codeberg"
	// Default: "github"
	Platform: config.PlatformGitHub,

	// Skip lists workflow names to skip during sync.
	// Example: []string{"lint", "test"}
	Skip: []string{},
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

// Sync regenerates CI workflows for the configured platform.
func Sync(ctx context.Context) error {
	return targets.SyncWorkflows(cfg)
}

// RunSerial runs all mutating targets serially.
func RunSerial(ctx context.Context) error {
	return targets.RunSerial(ctx, cfg, skip)
}

// RunParallel runs all non-mutating targets in parallel.
func RunParallel(ctx context.Context) error {
	return targets.RunParallel(ctx, cfg, skip)
}
