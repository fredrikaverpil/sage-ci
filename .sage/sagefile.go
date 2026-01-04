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

func main() {
	sg.GenerateMakefiles(
		sg.Makefile{
			Path:          sg.FromGitRoot("Makefile"),
			DefaultTarget: All,
		},
	)
}

func All(ctx context.Context) error {
	targets.RunSerial(ctx, cfg)
	targets.RunParallel(ctx, cfg)
	return targets.GitDiffCheck(ctx)
}

func GenerateWorkflows(ctx context.Context) error {
	return targets.GenerateWorkflows(cfg)
}

func UpdateSageCi(ctx context.Context) error {
	return targets.UpdateSageCi(ctx, cfg)
}

// Usage example: `make run target=go-format`
func Run(ctx context.Context, target string) error {
	return targets.Run(ctx, cfg, target)
}
