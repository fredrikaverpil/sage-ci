package workflows

import "fmt"

// Config configures the workflow generation.
type Config struct {
	// Ecosystem modules - explicit paths.
	// E.g. []string{".", "subdir/mylib"}
	GoModules []string `yaml:"go_modules"`
	// E.g. []string{"python", "tools/cli"}
	PythonModules []string `yaml:"python_modules"`
	// E.g. []string{"lua/plugin"}
	LuaModules []string `yaml:"lua_modules"`

	// Workflow selection (default: all enabled if empty).
	// E.g. []string{"stale", "release"}
	Skip []string `yaml:"skip"`

	// Options
	// default: ["stable"]
	GoVersions []string `yaml:"go_versions"`
	// default: ["3.12"]
	PythonVersions []string `yaml:"python_versions"`
	// default: ["ubuntu-latest"]
	OSVersions []string `yaml:"os_versions"`

	// OutputDir is the directory where workflows will be written.
	// default: ".github/workflows"
	OutputDir string `yaml:"output_dir"`
}

// Sync generates GitHub actions workflows based on the provided configuration.
func Sync(cfg Config) error {
	if cfg.OutputDir == "" {
		cfg.OutputDir = ".github/workflows"
	}
	if len(cfg.GoVersions) == 0 {
		cfg.GoVersions = []string{"stable"}
	}
	if len(cfg.PythonVersions) == 0 {
		cfg.PythonVersions = []string{"3.12"}
	}
	if len(cfg.OSVersions) == 0 {
		cfg.OSVersions = []string{"ubuntu-latest"}
	}

	if err := render(cfg); err != nil {
		return fmt.Errorf("render workflows: %w", err)
	}

	return nil
}
