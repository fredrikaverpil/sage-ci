package main

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/config"
	"github.com/fredrikaverpil/sage-ci/targets"
	"go.einride.tech/sage/sg"
)

var cfg = config.Config{
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

// Sync regenerates CI workflows for the configured platform.
func Sync(ctx context.Context) error {
	return targets.SyncWorkflows(cfg)
}

func RunSerial(ctx context.Context) error {
	return targets.RunSerial(ctx, cfg, skip)
}

func RunParallel(ctx context.Context) error {
	return targets.RunParallel(ctx, cfg, skip)
}
