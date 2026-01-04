package targets

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/config"
	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/tools/sguv"
)

func pythonSync(ctx context.Context, cfg config.Config, skip SkipTargets) error {
	for _, module := range cfg.PythonModules {
		if skip.ShouldSkip("PythonSync", module) {
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

func pythonFormat(ctx context.Context, cfg config.Config, skip SkipTargets) error {
	sg.Deps(ctx, func(ctx context.Context) error { return pythonSync(ctx, cfg, skip) })
	for _, module := range cfg.PythonModules {
		if skip.ShouldSkip("PythonFormat", module) {
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

func pythonLint(ctx context.Context, cfg config.Config, skip SkipTargets) error {
	sg.Deps(ctx, func(ctx context.Context) error { return pythonSync(ctx, cfg, skip) })
	for _, module := range cfg.PythonModules {
		if skip.ShouldSkip("PythonLint", module) {
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

func pythonMypy(ctx context.Context, cfg config.Config, skip SkipTargets) error {
	sg.Deps(ctx, func(ctx context.Context) error { return pythonSync(ctx, cfg, skip) })
	for _, module := range cfg.PythonModules {
		if skip.ShouldSkip("PythonMypy", module) {
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

func pythonTest(ctx context.Context, cfg config.Config, skip SkipTargets) error {
	sg.Deps(ctx, func(ctx context.Context) error { return pythonSync(ctx, cfg, skip) })
	for _, module := range cfg.PythonModules {
		if skip.ShouldSkip("PythonTest", module) {
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
