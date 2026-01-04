// Package sgtsqueryls provides a Sage tool for running ts_query_ls.
package sgtsqueryls

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/sgtool"
)

const name = "ts_query_ls"

// renovate: datasource=github-releases depName=ribru17/ts_query_ls
const version = "3.15.1"

func tsQueryLsPlatform(hostOS, hostArch string) (string, error) {
	var target string
	switch {
	case hostOS == "darwin" && hostArch == "arm64":
		target = "aarch64-apple-darwin"
	case hostOS == "darwin" && hostArch == "amd64":
		target = "x86_64-apple-darwin"
	case hostOS == "linux" && hostArch == "arm64":
		target = "aarch64-unknown-linux-gnu"
	case hostOS == "linux" && hostArch == "amd64":
		target = "x86_64-unknown-linux-gnu"
	case hostOS == "windows" && hostArch == "amd64":
		target = "x86_64-pc-windows-msvc"
	default:
		return "", fmt.Errorf("unsupported platform: %s/%s", hostOS, hostArch)
	}
	return target, nil
}

// Command returns an *exec.Cmd for ts_query_ls.
func Command(ctx context.Context, args ...string) *exec.Cmd {
	sg.Deps(ctx, PrepareCommand)
	return sg.Command(ctx, sg.FromBinDir(name), args...)
}

// PrepareCommand ensures ts_query_ls is installed.
func PrepareCommand(ctx context.Context) error {
	binDir := sg.FromToolsDir(name, version, "bin")
	binary := filepath.Join(binDir, name)
	hostOS := runtime.GOOS
	hostArch := runtime.GOARCH

	target, err := tsQueryLsPlatform(hostOS, hostArch)
	if err != nil {
		return err
	}

	// ts_query_ls-aarch64-apple-darwin.tar.gz, ts_query_ls-x86_64-pc-windows-msvc.zip
	var binURL string
	if hostOS == "windows" {
		binURL = fmt.Sprintf(
			"https://github.com/ribru17/ts_query_ls/releases/download/v%s/ts_query_ls-%s.zip",
			version, target,
		)
	} else {
		binURL = fmt.Sprintf(
			"https://github.com/ribru17/ts_query_ls/releases/download/v%s/ts_query_ls-%s.tar.gz",
			version, target,
		)
	}

	binaryName := name
	if hostOS == "windows" {
		binaryName = name + ".exe"
	}

	unarchive := sgtool.WithUntarGz()
	if hostOS == "windows" {
		unarchive = sgtool.WithUnzip()
	}

	if err := sgtool.FromRemote(
		ctx,
		binURL,
		sgtool.WithDestinationDir(binDir),
		unarchive,
		sgtool.WithRenameFile(binaryName, binaryName),
		sgtool.WithSkipIfFileExists(binary),
		sgtool.WithSymlink(binary),
	); err != nil {
		return fmt.Errorf("unable to download %s: %w", name, err)
	}
	return nil
}
