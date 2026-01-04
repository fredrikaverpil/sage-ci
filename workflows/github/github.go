// Package github generates GitHub Actions workflow files.
package github

import (
	"fmt"

	"github.com/fredrikaverpil/sage-ci/config"
)

const defaultOutputDir = ".github/workflows"

// outputDir can be overridden in tests.
var outputDir = defaultOutputDir

// Sync generates GitHub Actions workflows based on the provided configuration.
func Sync(cfg config.Config) error {
	cfg = cfg.WithDefaults()

	if err := render(cfg); err != nil {
		return fmt.Errorf("render github workflows: %w", err)
	}

	return nil
}
