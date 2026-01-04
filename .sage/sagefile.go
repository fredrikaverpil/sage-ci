package main

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/config"
	"github.com/fredrikaverpil/sage-ci/targets"
	"go.einride.tech/sage/sg"
)

var cfg = config.Config{
	GoModules:     []string{"."},
	SkipWorkflows: []string{"sage-ci-sync"},
}

func main() {
	sg.GenerateMakefiles(
		sg.Makefile{
			Path:          sg.FromGitRoot("Makefile"),
			DefaultTarget: All,
		},
	)
}

func All(ctx context.Context) error {
	if err := targets.RunSerial(ctx, cfg); err != nil {
		return err
	}
	if err := targets.RunParallel(ctx, cfg); err != nil {
		return err
	}
	sg.Deps(ctx, targets.GitDiffCheckTarget())
	return nil
}
