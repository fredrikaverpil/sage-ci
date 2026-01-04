// Package targets provides reusable CI/CD target functions for Sage-based projects.
package targets

import (
	"context"
	"fmt"
	"os"

	"github.com/fredrikaverpil/sage-ci/config"
	"github.com/fredrikaverpil/sage-ci/workflows/github"
	"go.einride.tech/sage/sg"
)

// --- Orchestration ---

// RunSerial runs all mutating targets serially for configured ecosystems.
func RunSerial(ctx context.Context, cfg config.Config) error {
	var deps []any
	if len(cfg.GoModules) > 0 {
		deps = append(deps,
			func(ctx context.Context) error { return GoModTidy(ctx, cfg) },
			func(ctx context.Context) error { return GoFormat(ctx, cfg) },
			func(ctx context.Context) error { return GoLint(ctx, cfg) },
		)
	}
	if len(cfg.PythonModules) > 0 {
		deps = append(deps,
			func(ctx context.Context) error { return PythonSync(ctx, cfg) },
			func(ctx context.Context) error { return PythonFormat(ctx, cfg) },
			func(ctx context.Context) error { return PythonLint(ctx, cfg) },
		)
	}
	if len(cfg.LuaModules) > 0 {
		deps = append(deps,
			func(ctx context.Context) error { return LuaFormat(ctx, cfg) },
		)
	}
	if len(deps) > 0 {
		sg.SerialDeps(ctx, deps...)
	}
	return nil
}

// RunParallel runs all non-mutating targets in parallel for configured ecosystems.
func RunParallel(ctx context.Context, cfg config.Config) error {
	var deps []any
	if len(cfg.GoModules) > 0 {
		deps = append(deps,
			func(ctx context.Context) error { return GoTest(ctx, cfg) },
			func(ctx context.Context) error { return GoVulncheck(ctx, cfg) },
		)
	}
	if len(cfg.PythonModules) > 0 {
		deps = append(deps,
			func(ctx context.Context) error { return PythonMypy(ctx, cfg) },
			func(ctx context.Context) error { return PythonTest(ctx, cfg) },
		)
	}
	if len(deps) > 0 {
		sg.Deps(ctx, deps...)
	}
	return nil
}

// --- Generate targets ---

// GenerateWorkflows generates CI workflows for the configured platforms.
// Defaults to GitHub if no platform is specified.
func GenerateWorkflows(cfg config.Config) error {
	cfg = cfg.WithDefaults()
	for _, platform := range cfg.Platforms {
		switch platform {
		case config.PlatformGitLab:
			return fmt.Errorf("gitlab workflows not yet implemented")
		case config.PlatformCodeberg:
			return fmt.Errorf("codeberg workflows not yet implemented")
		case config.PlatformGitHub:
			if err := github.Sync(cfg); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown platform: %s", platform)
		}
	}
	return nil
}

// --- Utility targets ---

// UpdateSageCi updates the sage-ci dependency, regenerates Makefiles and workflows.
func UpdateSageCi(ctx context.Context, cfg config.Config) error {
	// Skip dependency update if running from the sage-ci repo itself.
	if _, err := os.Stat(sg.FromGitRoot("cmd/sage-ci")); err == nil {
		sg.Logger(ctx).Println("skipping sage-ci dependency update (running from sage-ci repo)")
	} else {
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
	}

	sg.Logger(ctx).Println("generating targets.gen.go...")
	if err := GenerateTargetsFile(cfg, sg.FromGitRoot(".sage")); err != nil {
		return fmt.Errorf("generate targets file: %w", err)
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
