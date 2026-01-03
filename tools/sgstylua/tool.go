// Package sgstylua provides a Sage tool for running stylua.
package sgstylua

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

const name = "stylua"

// renovate: datasource=github-releases depName=JohnnyMorganz/StyLua
const version = "2.0.1"

const osWindows = "windows"

func styluaPlatform(hostOS, hostArch string) (string, string, error) {
	var osName string
	switch hostOS {
	case "darwin":
		osName = "macos"
	case "linux":
		osName = "linux"
	case osWindows:
		osName = osWindows
	default:
		return "", "", fmt.Errorf("unsupported OS: %s", hostOS)
	}

	var archName string
	switch hostArch {
	case "amd64":
		archName = "x86_64"
	case "arm64":
		archName = "aarch64"
	default:
		return "", "", fmt.Errorf("unsupported architecture: %s", hostArch)
	}

	return osName, archName, nil
}

// Command returns an *exec.Cmd for stylua.
func Command(ctx context.Context, args ...string) *exec.Cmd {
	sg.Deps(ctx, PrepareCommand)
	return sg.Command(ctx, sg.FromBinDir(name), args...)
}

// PrepareCommand ensures stylua is installed.
func PrepareCommand(ctx context.Context) error {
	binDir := sg.FromToolsDir(name, version, "bin")
	binary := filepath.Join(binDir, name)
	hostOS := runtime.GOOS
	hostArch := runtime.GOARCH

	// Map Go OS/arch to stylua naming convention
	osName, archName, err := styluaPlatform(hostOS, hostArch)
	if err != nil {
		return err
	}

	// stylua-macos-aarch64.zip, stylua-linux-x86_64.zip, etc.
	assetName := fmt.Sprintf("stylua-%s-%s.zip", osName, archName)
	binURL := fmt.Sprintf(
		"https://github.com/JohnnyMorganz/StyLua/releases/download/v%s/%s",
		version,
		assetName,
	)

	binaryName := name
	if hostOS == osWindows {
		binaryName = name + ".exe"
	}

	if err := sgtool.FromRemote(
		ctx,
		binURL,
		sgtool.WithDestinationDir(binDir),
		sgtool.WithUnzip(),
		sgtool.WithRenameFile(binaryName, binaryName),
		sgtool.WithSymlink(binary),
	); err != nil {
		return fmt.Errorf("unable to download %s: %w", name, err)
	}
	return nil
}

// Run runs stylua to check formatting in the current directory.
func Run(ctx context.Context) error {
	sg.Deps(ctx, PrepareCommand)
	cmd := Command(ctx, "--check", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// Fix runs stylua to fix formatting in the current directory.
func Fix(ctx context.Context) error {
	sg.Deps(ctx, PrepareCommand)
	cmd := Command(ctx, ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
