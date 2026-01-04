// Package targets provides reusable CI/CD target functions for Sage-based projects.
package targets

import (
	"context"
	"fmt"
	"os"

	"github.com/fredrikaverpil/sage-ci/workflows"
	"go.einride.tech/sage/sg"
)

// SkipTargets maps target names to modules that should be skipped.
// Key: Target name (e.g. "GoTest").
// Value: List of modules to skip. Use "*" to skip all modules.
type SkipTargets map[string][]string

// ShouldSkip returns true if the target should be skipped for the given module.
func (s SkipTargets) ShouldSkip(target, module string) bool {
	skippedModules, ok := s[target]
	if !ok {
		return false
	}
	for _, m := range skippedModules {
		if m == "*" || m == module {
			return true
		}
	}
	return false
}

// --- Orchestration ---

// RunSerial runs all mutating targets serially for configured ecosystems.
func RunSerial(ctx context.Context, cfg workflows.Config, skip SkipTargets) error {
	if len(cfg.GoModules) > 0 {
		sg.SerialDeps(ctx,
			func(ctx context.Context) error { return GoModTidy(ctx, cfg, skip) },
			func(ctx context.Context) error { return GoFormat(ctx, cfg, skip) },
			func(ctx context.Context) error { return GoLint(ctx, cfg, skip) },
		)
	}
	if len(cfg.PythonModules) > 0 {
		sg.SerialDeps(ctx,
			func(ctx context.Context) error { return PythonSync(ctx, cfg, skip) },
			func(ctx context.Context) error { return PythonFormat(ctx, cfg, skip) },
			func(ctx context.Context) error { return PythonLint(ctx, cfg, skip) },
		)
	}
	if len(cfg.LuaModules) > 0 {
		sg.SerialDeps(ctx,
			func(ctx context.Context) error { return LuaFormat(ctx, cfg, skip) },
		)
	}
	return nil
}

// RunParallel runs all non-mutating targets in parallel for configured ecosystems.
func RunParallel(ctx context.Context, cfg workflows.Config, skip SkipTargets) error {
	if len(cfg.GoModules) > 0 && len(cfg.PythonModules) > 0 {
		sg.Deps(ctx,
			func(ctx context.Context) error { return GoTest(ctx, cfg, skip) },
			func(ctx context.Context) error { return GoVulncheck(ctx, cfg, skip) },
			func(ctx context.Context) error { return PythonMypy(ctx, cfg, skip) },
			func(ctx context.Context) error { return PythonTest(ctx, cfg, skip) },
		)
	} else if len(cfg.GoModules) > 0 {
		sg.Deps(ctx,
			func(ctx context.Context) error { return GoTest(ctx, cfg, skip) },
			func(ctx context.Context) error { return GoVulncheck(ctx, cfg, skip) },
		)
	} else if len(cfg.PythonModules) > 0 {
		sg.Deps(ctx,
			func(ctx context.Context) error { return PythonMypy(ctx, cfg, skip) },
			func(ctx context.Context) error { return PythonTest(ctx, cfg, skip) },
		)
	}
	return nil
}

// --- Generate targets ---

// GenerateGHA generates GitHub Actions workflows.
func GenerateGHA(cfg workflows.Config) error {
	return workflows.Sync(cfg)
}

// --- Utility targets ---

// GitDiffCheck fails if there are uncommitted changes (only in CI).
func GitDiffCheck(ctx context.Context) error {
	hasDiff := sg.Command(ctx, "git", "diff", "--exit-code").Run() != nil ||
		sg.Command(ctx, "git", "diff", "--cached", "--exit-code").Run() != nil
	if !hasDiff {
		return nil
	}
	if os.Getenv("CI") == "" {
		sg.Logger(ctx).Println("warning: uncommitted changes detected")
		return nil
	}
	_ = sg.Command(ctx, "git", "diff").Run()
	return fmt.Errorf("uncommitted changes detected")
}
