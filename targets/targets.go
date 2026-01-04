// Package targets provides reusable CI/CD target functions for Sage-based projects.
package targets

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fredrikaverpil/sage-ci/config"
	"github.com/fredrikaverpil/sage-ci/workflows/github"
	"go.einride.tech/sage/sg"
)

// ErrUnknownTarget is returned when Run is called with an unrecognized target name.
var ErrUnknownTarget = fmt.Errorf("unknown target")

// Run executes a target by name. The target parameter uses kebab-case naming
// (e.g., "go-format", "python-lint") which maps to the corresponding function
// (e.g., GoFormat, PythonLint).
//
// Available targets:
//   - go-mod-tidy, go-format, go-lint, go-test, go-vulncheck
//   - python-sync, python-format, python-lint, python-mypy, python-test
//   - lua-format
//   - run-serial, run-parallel
func Run(ctx context.Context, cfg config.Config, skip SkipTargets, target string) error {
	switch strings.ToLower(target) {
	// Go targets.
	case "go-mod-tidy":
		return GoModTidy(ctx, cfg, skip)
	case "go-format":
		return GoFormat(ctx, cfg, skip)
	case "go-lint":
		return GoLint(ctx, cfg, skip)
	case "go-test":
		return GoTest(ctx, cfg, skip)
	case "go-vulncheck":
		return GoVulncheck(ctx, cfg, skip)
	// Python targets.
	case "python-sync":
		return PythonSync(ctx, cfg, skip)
	case "python-format":
		return PythonFormat(ctx, cfg, skip)
	case "python-lint":
		return PythonLint(ctx, cfg, skip)
	case "python-mypy":
		return PythonMypy(ctx, cfg, skip)
	case "python-test":
		return PythonTest(ctx, cfg, skip)
	// Lua targets.
	case "lua-format":
		return LuaFormat(ctx, cfg, skip)
	// Orchestration targets.
	case "run-serial":
		return RunSerial(ctx, cfg, skip)
	case "run-parallel":
		return RunParallel(ctx, cfg, skip)
	default:
		return fmt.Errorf("%w: %s", ErrUnknownTarget, target)
	}
}

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
func RunSerial(ctx context.Context, cfg config.Config, skip SkipTargets) error {
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
func RunParallel(ctx context.Context, cfg config.Config, skip SkipTargets) error {
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

// GenerateWorkflows generates CI workflows for the configured platform.
// Defaults to GitHub if no platform is specified.
func GenerateWorkflows(cfg config.Config) error {
	switch cfg.Platform {
	case config.PlatformGitLab:
		return fmt.Errorf("gitlab workflows not yet implemented")
	case config.PlatformCodeberg:
		return fmt.Errorf("codeberg workflows not yet implemented")
	case config.PlatformGitHub, "":
		return github.Sync(cfg)
	default:
		return fmt.Errorf("unknown platform: %s", cfg.Platform)
	}
}

// --- Utility targets ---

// UpdateSageCi updates the sage-ci dependency, regenerates Makefiles and workflows.
func UpdateSageCi(ctx context.Context, cfg config.Config) error {
	sg.Logger(ctx).Println("updating sage-ci dependency...")
	getCmd := sg.Command(ctx, "go", "get", "-u", "github.com/fredrikaverpil/sage-ci@latest")
	getCmd.Dir = sg.FromGitRoot(".sage")
	if err := getCmd.Run(); err != nil {
		return fmt.Errorf("update sage-ci dependency: %w", err)
	}

	sg.Logger(ctx).Println("running go mod tidy...")
	tidyCmd := sg.Command(ctx, "go", "mod", "tidy")
	tidyCmd.Dir = sg.FromGitRoot(".sage")
	if err := tidyCmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	sg.Logger(ctx).Println("regenerating Makefile(s)...")
	makefileCmd := sg.Command(ctx, "go", "run", ".")
	makefileCmd.Dir = sg.FromGitRoot(".sage")
	if err := makefileCmd.Run(); err != nil {
		return fmt.Errorf("regenerate makefiles: %w", err)
	}

	sg.Logger(ctx).Println("regenerating workflows...")
	if err := GenerateWorkflows(cfg); err != nil {
		return fmt.Errorf("regenerate workflows: %w", err)
	}

	return nil
}

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
