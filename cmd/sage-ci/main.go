package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fredrikaverpil/sage-ci/workflows"
	"gopkg.in/yaml.v3"
)

//go:embed templates/sagefile.go
var sagefileContent string

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		if err := runInit(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "sync":
		if err := runSync(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`Usage: sage-ci <command> [flags]

Commands:
  init    Bootstrap a new project with .sage/ directory
  sync    Sync workflows to .github/workflows`)
}

func runInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := os.MkdirAll(".sage", 0o755); err != nil {
		return fmt.Errorf("create .sage directory: %w", err)
	}

	// Check if sagefile.go already exists
	if _, err := os.Stat(".sage/sagefile.go"); err == nil {
		return errors.New(".sage/sagefile.go already exists")
	}

	if err := os.WriteFile(".sage/sagefile.go", []byte(sagefileContent), 0o644); err != nil {
		return fmt.Errorf("write .sage/sagefile.go: %w", err)
	}

	fmt.Println("Initialized .sage/sagefile.go")
	fmt.Println("Edit .sage/sagefile.go to configure your project, then run:")
	fmt.Println("  go run ./.sage")
	return nil
}

func runSync(args []string) error {
	fs := flag.NewFlagSet("sync", flag.ExitOnError)
	var (
		configFile    = fs.String("config", "", "Path to configuration file (e.g. .sage-ci.yml)")
		goModules     = fs.String("go-modules", "", "Comma-separated list of Go modules")
		pythonModules = fs.String("python-modules", "", "Comma-separated list of Python modules")
		luaModules    = fs.String("lua-modules", "", "Comma-separated list of Lua modules")
		skip          = fs.String("skip", "", "Comma-separated list of workflows to skip")
		outputDir     = fs.String("output-dir", "", "Output directory for workflows")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	var cfg workflows.Config

	// Load from config file if provided
	if *configFile != "" {
		data, err := os.ReadFile(*configFile)
		if err != nil {
			return fmt.Errorf("read config file: %w", err)
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("parse config file: %w", err)
		}
	}

	// Override with flags if provided
	if *goModules != "" {
		cfg.GoModules = strings.Split(*goModules, ",")
	}
	if *pythonModules != "" {
		cfg.PythonModules = strings.Split(*pythonModules, ",")
	}
	if *luaModules != "" {
		cfg.LuaModules = strings.Split(*luaModules, ",")
	}
	if *skip != "" {
		cfg.Skip = strings.Split(*skip, ",")
	}
	if *outputDir != "" {
		cfg.OutputDir = *outputDir
	}

	return workflows.Sync(cfg)
}
