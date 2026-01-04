package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/fredrikaverpil/sage-ci/config"
)

type templateData struct {
	// Metadata
	GeneratedBy string
	Timestamp   string

	// Module paths
	GoModules     []string
	PythonModules []string
	LuaModules    []string

	// Version matrices
	GoVersions     []string
	PythonVersions []string
	OSVersions     []string

	// Skipped targets (fully skipped for all modules)
	SkipGoTest       bool
	SkipGoLint       bool
	SkipGoFormat     bool
	SkipGoVulncheck  bool
	SkipPythonTest   bool
	SkipPythonLint   bool
	SkipPythonFormat bool
	SkipPythonMypy   bool
	SkipLuaFormat    bool
}

func render(cfg config.Config) error {
	data := templateData{
		GeneratedBy:    "sage-ci",
		Timestamp:      time.Now().Format(time.RFC3339),
		GoModules:      cfg.GoModules,
		PythonModules:  cfg.PythonModules,
		LuaModules:     cfg.LuaModules,
		GoVersions:     cfg.GoVersions,
		PythonVersions: cfg.PythonVersions,
		OSVersions:     cfg.OSVersions,

		// Check if targets are fully skipped
		SkipGoTest:       cfg.SkipTargets.IsFullySkipped("GoTest", cfg.GoModules),
		SkipGoLint:       cfg.SkipTargets.IsFullySkipped("GoLint", cfg.GoModules),
		SkipGoFormat:     cfg.SkipTargets.IsFullySkipped("GoFormat", cfg.GoModules),
		SkipGoVulncheck:  cfg.SkipTargets.IsFullySkipped("GoVulncheck", cfg.GoModules),
		SkipPythonTest:   cfg.SkipTargets.IsFullySkipped("PythonTest", cfg.PythonModules),
		SkipPythonLint:   cfg.SkipTargets.IsFullySkipped("PythonLint", cfg.PythonModules),
		SkipPythonFormat: cfg.SkipTargets.IsFullySkipped("PythonFormat", cfg.PythonModules),
		SkipPythonMypy:   cfg.SkipTargets.IsFullySkipped("PythonMypy", cfg.PythonModules),
		SkipLuaFormat:    cfg.SkipTargets.IsFullySkipped("LuaFormat", cfg.LuaModules),
	}

	funcMap := template.FuncMap{
		"toJSON": func(v any) (string, error) {
			b, err := json.Marshal(v)
			if err != nil {
				return "", fmt.Errorf("marshal to JSON: %w", err)
			}
			return string(b), nil
		},
	}

	// Walk the templates directory
	err := fs.WalkDir(templatesFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Skip hidden files or non-templates
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		// Read template content
		tmplContent, err := templatesFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read template %s: %w", path, err)
		}

		// Parse template
		t, err := template.New(filepath.Base(path)).Funcs(funcMap).Parse(string(tmplContent))
		if err != nil {
			return fmt.Errorf("parse template %s: %w", path, err)
		}

		// Render
		var buf bytes.Buffer
		if err := t.Execute(&buf, data); err != nil {
			return fmt.Errorf("execute template %s: %w", path, err)
		}

		relPath, err := filepath.Rel("templates", path)
		if err != nil {
			return fmt.Errorf("get relative path for %s: %w", path, err)
		}
		parts := strings.Split(relPath, string(os.PathSeparator))

		// Determine output filename:
		// - generic/*.yml.tmpl -> sage-ci-*.yml
		// - <ecosystem>/*.yml.tmpl -> sage-ci-<ecosystem>-*.yml
		var fileName string
		if len(parts) == 2 {
			category := parts[0]
			name := strings.TrimSuffix(parts[1], ".tmpl")

			if category == "generic" {
				fileName = fmt.Sprintf("sage-ci-%s", name)
			} else {
				fileName = fmt.Sprintf("sage-ci-%s-%s", category, name)
			}
		} else {
			// Fallback
			fileName = "sage-ci-" + strings.TrimSuffix(filepath.Base(path), ".tmpl")
		}

		// Check for skip
		baseName := strings.TrimSuffix(fileName, ".yml")
		if slices.Contains(cfg.SkipWorkflows, baseName) {
			return nil
		}

		// Skip ecosystem-specific workflows if no modules configured
		if (parts[0] == "go" && len(cfg.GoModules) == 0) ||
			(parts[0] == "python" && len(cfg.PythonModules) == 0) ||
			(parts[0] == "lua" && len(cfg.LuaModules) == 0) {
			return nil
		}

		outputPath := filepath.Join(outputDir, fileName)

		// Ensure output dir exists
		if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}

		if err := os.WriteFile(outputPath, buf.Bytes(), 0o644); err != nil {
			return fmt.Errorf("write workflow %s: %w", outputPath, err)
		}

		return nil
	})

	return err
}
