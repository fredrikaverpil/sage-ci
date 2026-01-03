// Package sgmypy provides a Sage tool for running mypy via uv.
package sgmypy

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/tools/sguv"
)

const name = "mypy"

// renovate: datasource=pypi depName=mypy
const version = "1.19.0"

// Command returns an *exec.Cmd for mypy via uv tool run.
func Command(ctx context.Context, args ...string) *exec.Cmd {
	sg.Deps(ctx, sguv.PrepareCommand)
	return sguv.Command(ctx, append([]string{"tool", "run", fmt.Sprintf("%s@%s", name, version)}, args...)...)
}

// Run runs mypy on the current directory.
func Run(ctx context.Context) error {
	cmd := Command(ctx, ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// Check runs mypy with strict mode on the current directory.
func Check(ctx context.Context) error {
	cmd := Command(ctx, "--strict", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
