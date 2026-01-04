package main

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/targets"
	"github.com/fredrikaverpil/sage-ci/workflows"
	"go.einride.tech/sage/sg"
)

var cfg = workflows.Config{
	GoModules: []string{"."},
}

var skip = targets.SkipTargets{}

func main() {
	sg.GenerateMakefiles(
		sg.Makefile{
			Path:          sg.FromGitRoot("Makefile"),
			DefaultTarget: All,
		},
	)
}

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

func RunSerial(ctx context.Context) error {
	return targets.RunSerial(ctx, cfg, skip)
}

func RunParallel(ctx context.Context) error {
	return targets.RunParallel(ctx, cfg, skip)
}
