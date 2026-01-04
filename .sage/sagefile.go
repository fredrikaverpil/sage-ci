package main

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/workflows"
	"go.einride.tech/sage/sg"
)

var cfg = workflows.Config{
	GoModules: []string{"."},
}

var skipTargets = map[string][]string{}

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
	sg.Deps(ctx, RunSyncedSerial)
	sg.Deps(ctx, RunSynced)
	return GitDiffCheck(ctx)
}
