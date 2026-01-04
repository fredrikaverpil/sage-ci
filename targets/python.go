package targets

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/config"
	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/tools/sguv"
)

// PythonSync runs uv sync for all configured Python modules.
func PythonSync(ctx context.Context, cfg config.Config) error {
	for _, module := range cfg.PythonModules {
		if cfg.SkipTargets.ShouldSkip("PythonSync", module) {
			continue
		}
		sg.Logger(ctx).Printf("running uv sync in %s...", module)
		cmd := sguv.Command(ctx, "sync", "--all-groups")
		cmd.Dir = sg.FromGitRoot(module)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// PythonFormat runs ruff format for all configured Python modules.
func PythonFormat(ctx context.Context, cfg config.Config) error {
	sg.Deps(ctx, func(ctx context.Context) error { return PythonSync(ctx, cfg) })
	for _, module := range cfg.PythonModules {
		if cfg.SkipTargets.ShouldSkip("PythonFormat", module) {
			continue
		}
		sg.Logger(ctx).Printf("applying ruff format in %s...", module)
		cmd := sguv.Command(ctx, "run", "ruff", "format", ".")
		cmd.Dir = sg.FromGitRoot(module)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// PythonLint runs ruff check for all configured Python modules.
func PythonLint(ctx context.Context, cfg config.Config) error {
	sg.Deps(ctx, func(ctx context.Context) error { return PythonSync(ctx, cfg) })
	for _, module := range cfg.PythonModules {
		if cfg.SkipTargets.ShouldSkip("PythonLint", module) {
			continue
		}
		sg.Logger(ctx).Printf("running ruff check --fix in %s...", module)
		cmd := sguv.Command(ctx, "run", "ruff", "check", "--fix", ".")
		cmd.Dir = sg.FromGitRoot(module)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// PythonMypy runs mypy for all configured Python modules.
func PythonMypy(ctx context.Context, cfg config.Config) error {
	sg.Deps(ctx, func(ctx context.Context) error { return PythonSync(ctx, cfg) })
	for _, module := range cfg.PythonModules {
		if cfg.SkipTargets.ShouldSkip("PythonMypy", module) {
			continue
		}
		sg.Logger(ctx).Printf("running mypy in %s...", module)
		cmd := sguv.Command(ctx, "run", "mypy", ".")
		cmd.Dir = sg.FromGitRoot(module)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// PythonTest runs pytest for all configured Python modules.
func PythonTest(ctx context.Context, cfg config.Config) error {
	sg.Deps(ctx, func(ctx context.Context) error { return PythonSync(ctx, cfg) })
	for _, module := range cfg.PythonModules {
		if cfg.SkipTargets.ShouldSkip("PythonTest", module) {
			continue
		}
		sg.Logger(ctx).Printf("running pytest in %s...", module)
		cmd := sguv.Command(ctx, "run", "pytest", "-v")
		cmd.Dir = sg.FromGitRoot(module)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
