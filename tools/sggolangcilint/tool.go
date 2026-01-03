// Package sggolangcilint provides a Sage tool for running golangci-lint v2.
package sggolangcilint

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/sgtool"
)

const (
	name      = "golangci-lint"
	osWindows = "windows"
)

// renovate: datasource=github-releases depName=golangci/golangci-lint
const version = "2.7.1"

// Command returns an *exec.Cmd for golangci-lint.
func Command(ctx context.Context, args ...string) *exec.Cmd {
	sg.Deps(ctx, PrepareCommand)
	return sg.Command(ctx, sg.FromBinDir(name), args...)
}

// PrepareCommand ensures golangci-lint is installed.
func PrepareCommand(ctx context.Context) error {
	binDir := sg.FromToolsDir(name, version, "bin")
	binary := filepath.Join(binDir, name)
	hostOS := runtime.GOOS
	hostArch := runtime.GOARCH

	// Windows uses .zip, others use .tar.gz
	var binURL string
	if hostOS == osWindows {
		binURL = fmt.Sprintf(
			"https://github.com/golangci/golangci-lint/releases/download/v%s/golangci-lint-%s-%s-%s.zip",
			version, version, hostOS, hostArch,
		)
	} else {
		binURL = fmt.Sprintf(
			"https://github.com/golangci/golangci-lint/releases/download/v%s/golangci-lint-%s-%s-%s.tar.gz",
			version, version, hostOS, hostArch,
		)
	}

	archiveDir := fmt.Sprintf("golangci-lint-%s-%s-%s", version, hostOS, hostArch)
	unarchive := sgtool.WithUntarGz()
	if hostOS == osWindows {
		unarchive = sgtool.WithUnzip()
	}

	if err := sgtool.FromRemote(
		ctx,
		binURL,
		sgtool.WithDestinationDir(binDir),
		unarchive,
		sgtool.WithRenameFile(fmt.Sprintf("%s/golangci-lint", archiveDir), name),
		sgtool.WithSkipIfFileExists(binary),
		sgtool.WithSymlink(binary),
	); err != nil {
		return fmt.Errorf("unable to download %s: %w", name, err)
	}
	return nil
}

// Run runs golangci-lint in the current directory.
func Run(ctx context.Context) error {
	sg.Deps(ctx, PrepareCommand)
	cmd := Command(ctx, "run", "--allow-parallel-runners", "-c", configPath(), "./...")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// Fix runs golangci-lint with --fix in the current directory.
func Fix(ctx context.Context) error {
	sg.Deps(ctx, PrepareCommand)
	cmd := Command(ctx, "run", "--fix", "--allow-parallel-runners", "-c", configPath(), "./...")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func configPath() string {
	path := sg.FromGitRoot(".golangci.yml")
	if _, err := os.Stat(path); err == nil {
		return path
	}
	return sg.FromGitRoot("tools/sggolangcilint/golangci.yml")
}

// Fmt runs golangci-lint fmt in the current directory.
func Fmt(ctx context.Context) error {
	sg.Deps(ctx, PrepareCommand)
	cmd := Command(ctx, "fmt", "./...")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
