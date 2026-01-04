// Package config provides shared configuration for sage-ci.
package config

// Platform represents a CI/CD platform for workflow generation.
type Platform string

const (
	// PlatformGitHub generates GitHub Actions workflows.
	PlatformGitHub Platform = "github"
	// PlatformGitLab generates GitLab CI workflows.
	PlatformGitLab Platform = "gitlab"
	// PlatformCodeberg generates Codeberg/Woodpecker workflows.
	PlatformCodeberg Platform = "codeberg"
)

// Config configures sage-ci targets and workflow generation.
type Config struct {
	// Ecosystem modules - explicit paths.
	// E.g. []string{".", "subdir/mylib"}
	GoModules []string
	// E.g. []string{"python", "tools/cli"}
	PythonModules []string
	// E.g. []string{"lua/plugin"}
	LuaModules []string

	// Workflow platform to generate for.
	// Default: "github"
	Platform []Platform

	// Workflow selection (default: all enabled if empty).
	// E.g. []string{"sage-ci-stale", "sage-ci-release"}
	// You can also use a string, and if found, the workflow will be skipped.
	SkipWorkflows []string

	// Options
	// default: ["stable"]
	GoVersions []string
	// default: ["3.12"]
	PythonVersions []string
	// default: ["ubuntu-latest"]
	OSVersions []string

	// OutputDir is the directory where workflows will be written.
	// Platform-specific defaults apply if empty.
	OutputDir string
}

// WithDefaults returns a copy of the config with default values applied.
func (c Config) WithDefaults() Config {
	if len(c.GoVersions) == 0 {
		c.GoVersions = []string{"stable"}
	}
	if len(c.PythonVersions) == 0 {
		c.PythonVersions = []string{"3.14"}
	}
	if len(c.OSVersions) == 0 {
		c.OSVersions = []string{"ubuntu-latest"}
	}
	if len(c.Platform) == 0 {
		c.Platform = []Platform{PlatformGitHub}
	}
	return c
}

// HasGo returns true if Go modules are configured.
func (c Config) HasGo() bool {
	return len(c.GoModules) > 0
}

// HasPython returns true if Python modules are configured.
func (c Config) HasPython() bool {
	return len(c.PythonModules) > 0
}

// HasLua returns true if Lua modules are configured.
func (c Config) HasLua() bool {
	return len(c.LuaModules) > 0
}
