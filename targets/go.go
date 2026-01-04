package targets

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/tools/sggolangcilint"
	"github.com/fredrikaverpil/sage-ci/workflows"
	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/tools/sggo"
)

// GoModTidy runs go mod tidy.
func GoModTidy(ctx context.Context, cfg workflows.Config, skip SkipTargets) error {
	for _, module := range cfg.GoModules {
		if skip.ShouldSkip("GoModTidy", module) {
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

// GoLint runs golangci-lint with --fix.
func GoLint(ctx context.Context, cfg workflows.Config, skip SkipTargets) error {
	for _, module := range cfg.GoModules {
		if skip.ShouldSkip("GoLint", module) {
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

// GoFormat applies Go formatting using gofmt.
func GoFormat(ctx context.Context, cfg workflows.Config, skip SkipTargets) error {
	for _, module := range cfg.GoModules {
		if skip.ShouldSkip("GoFormat", module) {
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

// GoTest runs Go tests.
func GoTest(ctx context.Context, cfg workflows.Config, skip SkipTargets) error {
	for _, module := range cfg.GoModules {
		if skip.ShouldSkip("GoTest", module) {
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

// GoVulncheck runs govulncheck.
func GoVulncheck(ctx context.Context, cfg workflows.Config, skip SkipTargets) error {
	for _, module := range cfg.GoModules {
		if skip.ShouldSkip("GoVulncheck", module) {
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
