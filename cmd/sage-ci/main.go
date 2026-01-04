package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

//go:embed templates/sagefile.go.tmpl
var sagefileContent string

//go:embed templates/go.mod.tmpl
var gomodContent string

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
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`Usage: sage-ci <command> [flags]

Commands:
  init    Bootstrap a new project with .sage/ directory`)
}

func runInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := os.MkdirAll(".sage", 0o755); err != nil {
		return fmt.Errorf("create .sage directory: %w", err)
	}

	// Check if sagefile.go already exists.
	if _, err := os.Stat(".sage/sagefile.go"); err == nil {
		return errors.New(".sage/sagefile.go already exists")
	}

	if err := os.WriteFile(".sage/go.mod", []byte(gomodContent), 0o644); err != nil {
		return fmt.Errorf("write .sage/go.mod: %w", err)
	}

	if err := os.WriteFile(".sage/sagefile.go", []byte(sagefileContent), 0o644); err != nil {
		return fmt.Errorf("write .sage/sagefile.go: %w", err)
	}

	fmt.Println("Initialized .sage/go.mod and .sage/sagefile.go")

	fmt.Println("Running go get -u ./... in .sage/...")
	getCmd := exec.Command("go", "get", "-u", "./...")
	getCmd.Dir = ".sage"
	getCmd.Stdout = os.Stdout
	getCmd.Stderr = os.Stderr
	if err := getCmd.Run(); err != nil {
		return fmt.Errorf("go get -u: %w", err)
	}

	fmt.Println("Running go mod tidy in .sage/...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = ".sage"
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr
	if err := tidyCmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	fmt.Println("Edit .sage/sagefile.go to configure your project, then run:")
	fmt.Println("  go run ./.sage")
	return nil
}
