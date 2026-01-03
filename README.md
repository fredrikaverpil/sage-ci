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

This creates:

- `.sage/sagefile.go` - Your project configuration and targets (edit this)
- `.sage/synced.gen.go` - Generated utility functions (do not edit)

### 3. Configure your project

Edit `.sage/sagefile.go` to specify your modules and customize targets:

```go
var cfg = workflows.Config{
    GoModules: []string{"."},
    // PythonModules: []string{"python"},
    // LuaModules:    []string{"lua"},
}

func All(ctx context.Context) error {
    sg.Deps(ctx, SyncGHA, GoTest)
    sg.SerialDeps(ctx, GoModTidy)
    return nil
}
```

You can add custom targets to `sagefile.go` or create additional `.go` files in
`.sage/`.

### 4. Generate Makefile and sync workflows

```bash
# Generate the Makefile
go run ./.sage

# Sync workflows (updates synced.gen.go and generates .github/workflows/)
sage-ci sync
```

## Updating

To update `synced.gen.go` with the latest utility functions:

```bash
sage-ci sync
```

This updates `.sage/synced.gen.go` while preserving your `.sage/sagefile.go`.

## CLI Reference

```
sage-ci init    # Bootstrap .sage/ directory
sage-ci sync    # Update synced.gen.go and sync workflows to .github/workflows
```

The `sync` command also accepts flags for standalone use without a `.sage/`
directory:

```bash
sage-ci sync --go-modules=. --python-modules=scripts
sage-ci sync --config=.sage-ci.yml
```
