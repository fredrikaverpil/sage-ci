// Package github generates GitHub Actions workflow files.
package github

import (
	"fmt"

	"github.com/fredrikaverpil/sage-ci/config"
)

const defaultOutputDir = ".github/workflows"

// Sync generates GitHub Actions workflows based on the provided configuration.
func Sync(cfg config.Config) error {
	if cfg.OutputDir == "" {
		cfg.OutputDir = defaultOutputDir
	}
	cfg = cfg.WithDefaults()

	if err := render(cfg); err != nil {
		return fmt.Errorf("render github workflows: %w", err)
	}

	return nil
}
