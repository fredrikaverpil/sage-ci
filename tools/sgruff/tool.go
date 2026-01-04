// Package sgruff provides a Sage tool for running ruff via uv.
package sgruff

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/tools/sguv"
)

const name = "ruff"

// renovate: datasource=pypi depName=ruff
const version = "0.14.10"

// Command returns an *exec.Cmd for ruff via uv tool run.
func Command(ctx context.Context, args ...string) *exec.Cmd {
	sg.Deps(ctx, sguv.PrepareCommand)
	return sguv.Command(ctx, append([]string{"tool", "run", fmt.Sprintf("%s@%s", name, version)}, args...)...)
}

// Check runs ruff check on the current directory.
func Check(ctx context.Context) error {
	cmd := Command(ctx, "check", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// Fix runs ruff check --fix on the current directory.
func Fix(ctx context.Context) error {
	cmd := Command(ctx, "check", "--fix", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// Format runs ruff format on the current directory.
func Format(ctx context.Context) error {
	cmd := Command(ctx, "format", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// FormatCheck runs ruff format --check on the current directory.
func FormatCheck(ctx context.Context) error {
	cmd := Command(ctx, "format", "--check", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
