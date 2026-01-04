// Package sgtreesittercli provides a Sage tool for running tree-sitter CLI.
package sgtreesittercli

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"go.einride.tech/sage/sg"
)

const name = "tree-sitter"

// renovate: datasource=github-releases depName=tree-sitter/tree-sitter
const version = "0.26.3"

func treeSitterPlatform(hostOS, hostArch string) (string, error) {
	var platform string
	switch {
	case hostOS == "darwin" && hostArch == "arm64":
		platform = "macos-arm64"
	case hostOS == "darwin" && hostArch == "amd64":
		platform = "macos-x64"
	case hostOS == "linux" && hostArch == "arm64":
		platform = "linux-arm64"
	case hostOS == "linux" && hostArch == "amd64":
		platform = "linux-x64"
	case hostOS == "windows" && hostArch == "arm64":
		platform = "windows-arm64"
	case hostOS == "windows" && hostArch == "amd64":
		platform = "windows-x64"
	default:
		return "", fmt.Errorf("unsupported platform: %s/%s", hostOS, hostArch)
	}
	return platform, nil
}

// Command returns an *exec.Cmd for tree-sitter.
func Command(ctx context.Context, args ...string) *exec.Cmd {
	sg.Deps(ctx, PrepareCommand)
	return sg.Command(ctx, sg.FromBinDir(name), args...)
}

// PrepareCommand ensures tree-sitter is installed.
func PrepareCommand(ctx context.Context) error {
	binDir := sg.FromToolsDir(name, version, "bin")
	hostOS := runtime.GOOS
	hostArch := runtime.GOARCH

	binaryName := name
	if hostOS == "windows" {
		binaryName = name + ".exe"
	}
	binary := filepath.Join(binDir, binaryName)

	// Skip if already installed.
	if _, err := os.Stat(binary); err == nil {
		// Ensure symlink exists.
		symlink := sg.FromBinDir(name)
		if _, err := os.Lstat(symlink); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(symlink), 0o755); err != nil {
				return err
			}
			if err := os.Symlink(binary, symlink); err != nil {
				return err
			}
		}
		return nil
	}

	platform, err := treeSitterPlatform(hostOS, hostArch)
	if err != nil {
		return err
	}

	// tree-sitter-macos-arm64.gz, tree-sitter-linux-x64.gz, etc.
	binURL := fmt.Sprintf(
		"https://github.com/tree-sitter/tree-sitter/releases/download/v%s/tree-sitter-%s.gz",
		version, platform,
	)

	if err := downloadGzipBinary(binURL, binDir, binaryName); err != nil {
		return fmt.Errorf("unable to download %s: %w", name, err)
	}

	// Create symlink.
	symlink := sg.FromBinDir(name)
	if err := os.MkdirAll(filepath.Dir(symlink), 0o755); err != nil {
		return err
	}
	_ = os.Remove(symlink)
	if err := os.Symlink(binary, symlink); err != nil {
		return err
	}

	return nil
}

func downloadGzipBinary(url, destDir, binaryName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gr.Close()

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}

	destPath := filepath.Join(destDir, binaryName)
	out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, gr); err != nil {
		return err
	}

	return nil
}
