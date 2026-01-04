package targets

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/config"
	"github.com/fredrikaverpil/sage-ci/tools/sggolangcilint"
	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/tools/sggo"
)

// GoModTidy runs go mod tidy for all configured Go modules.
func GoModTidy(ctx context.Context, cfg config.Config) error {
	for _, module := range cfg.GoModules {
		if cfg.SkipTargets.ShouldSkip("GoModTidy", module) {
			continue
		}
		sg.Logger(ctx).Printf("running go mod tidy in %s...", module)
		cmd := sg.Command(ctx, "go", "mod", "tidy", "-v")
		cmd.Dir = sg.FromGitRoot(module)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// GoLint runs golangci-lint for all configured Go modules.
func GoLint(ctx context.Context, cfg config.Config) error {
	for _, module := range cfg.GoModules {
		if cfg.SkipTargets.ShouldSkip("GoLint", module) {
			continue
		}
		sg.Logger(ctx).Printf("running golangci-lint --fix in %s...", module)
		cmd := sggolangcilint.Command(ctx, "run", "--fix", "--allow-parallel-runners", "./...")
		cmd.Dir = sg.FromGitRoot(module)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// GoFormat runs gofmt for all configured Go modules.
func GoFormat(ctx context.Context, cfg config.Config) error {
	for _, module := range cfg.GoModules {
		if cfg.SkipTargets.ShouldSkip("GoFormat", module) {
			continue
		}
		sg.Logger(ctx).Printf("applying gofmt in %s...", module)
		cmd := sg.Command(ctx, "gofmt", "-w", ".")
		cmd.Dir = sg.FromGitRoot(module)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// GoTest runs go test for all configured Go modules.
func GoTest(ctx context.Context, cfg config.Config) error {
	for _, module := range cfg.GoModules {
		if cfg.SkipTargets.ShouldSkip("GoTest", module) {
			continue
		}
		sg.Logger(ctx).Printf("running go test in %s...", module)
		cmd := sggo.TestCommand(ctx)
		cmd.Dir = sg.FromGitRoot(module)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// GoVulncheck runs govulncheck for all configured Go modules.
func GoVulncheck(ctx context.Context, cfg config.Config) error {
	for _, module := range cfg.GoModules {
		if cfg.SkipTargets.ShouldSkip("GoVulncheck", module) {
			continue
		}
		sg.Logger(ctx).Printf("running govulncheck in %s...", module)
		cmd := sg.Command(ctx, "go", "run", "golang.org/x/vuln/cmd/govulncheck@latest", "./...")
		cmd.Dir = sg.FromGitRoot(module)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
