package main

import (
	"context"

	"go.einride.tech/sage/sg"
)

// GoBuild builds the sage-ci tool.
func BuildSageCI(ctx context.Context) error {
	sg.Logger(ctx).Printf("building sage-ci...")
	cmd := sg.Command(ctx, "go", "build", "-o", "bin/sage-ci", "./cmd/sage-ci")
	cmd.Dir = sg.FromGitRoot()
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
