# sage-ci

This repo holds tools and configurations which I want synced across several
projects.

The engine powering everything is [Sage](https://github.com/einride/sage).

## Renovate tool updates

Each tool lives in `tools/<toolname>/tool.go` and follows this pattern:

```go
// renovate: datasource=github-releases depName=owner/repo
const version = "1.2.3"
```

Renovate will automatically create PRs when new versions are available.

When adding new tools, use the appropriate
[datasource](https://docs.renovatebot.com/modules/datasource/) for version
lookups.

## Bootstrap

To bootstrap `sage-ci` into a new repository, use the `sage-ci` CLI.

### 1. Install or run the CLI

```bash
# Run directly via its remote path
go run github.com/fredrikaverpil/sage-ci/cmd/sage-ci@latest <command>

# Or build and run locally (if you have the repo cloned)
go build -o sage-ci ./cmd/sage-ci
```

### 2. Initialize the project

Run the `init` command to create the `.sage/` directory:

```bash
sage-ci init
```

This creates `.sage/sagefile.go` with your project configuration.

### 3. Configure your project

Edit `.sage/sagefile.go` to specify your modules and customize targets:

```go
import (
    "github.com/fredrikaverpil/sage-ci/targets"
    "github.com/fredrikaverpil/sage-ci/workflows"
)

var cfg = workflows.Config{
    GoModules: []string{"."},
    // PythonModules: []string{"python"},
    // LuaModules:    []string{"lua"},
}

var skip = targets.SkipTargets{}

func All(ctx context.Context) error {
    sg.Deps(ctx, RunSerial)
    sg.Deps(ctx, RunParallel)
    return targets.GitDiffCheck(ctx)
}

func RunSerial(ctx context.Context) error {
    return targets.RunSerial(ctx, cfg, skip)
}

func RunParallel(ctx context.Context) error {
    return targets.RunParallel(ctx, cfg, skip)
}
```

You can add custom targets to `sagefile.go` or create additional `.go` files in
`.sage/`.

### 4. Generate Makefile and sync workflows

```bash
# Generate the Makefile
go run ./.sage

# Sync workflows to .github/workflows/
sage-ci sync
```

## Updating

Update your `sage-ci` dependency to get the latest targets:

```bash
cd .sage && go get -u github.com/fredrikaverpil/sage-ci@latest
```

Then regenerate workflows:

```bash
sage-ci sync
```

## CLI Reference

```
sage-ci init    # Bootstrap .sage/ directory
sage-ci sync    # Sync workflows to .github/workflows
```

The `sync` command also accepts flags for standalone use without a `.sage/`
directory:

```bash
sage-ci sync --go-modules=. --python-modules=scripts
sage-ci sync --config=.sage-ci.yml
```
