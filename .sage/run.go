package main

import (
	"context"

	"go.einride.tech/sage/sg"
)

// Sync runs sage-ci sync to regenerate workflows.
func Sync(ctx context.Context) error {
	sg.Logger(ctx).Printf("running sage-ci sync...")
	cmd := sg.Command(ctx, "go", "run", "./cmd/sage-ci", "sync")
	cmd.Dir = sg.FromGitRoot()
	return cmd.Run()
}
