# sage-ci (powered by [Sage](https://github.com/einride/sage) 🌿)

This repo holds opinionated tools and workflows which I want synced across
several projects. Features include:

- Opt-in ecosystem support for Go, Lua, Python and CI workflows
- Project Makefile for running commands in the project
- Sage-powered tools declaration
- CI workflow templates

## Quickstart

### Bootstrap to your project

```bash
go run github.com/fredrikaverpil/sage-ci/cmd/sage-ci@latest init
```

This creates `.sage/go.mod` and `.sage/sagefile.go` with your project
configuration.

### Configure

Edit `.sage/sagefile.go` to specify your modules and customize targets:

```go
var cfg = config.Config{
    GoModules:     []string{"."},
    PythonModules: []string{"tests"},
    LuaModules:    []string{"lua"},
    Platforms:     []config.Platform{config.PlatformGitHub},
}
```

See [config/config.go](config/config.go) for all configuration options.

You can add custom targets to `sagefile.go` or create additional `.go` files in
`.sage/`. Sage-ci provides opinionated targets in `RunSerial` and `RunParallel`.

### Generate Makefile, targets and workflows

```bash
# Generate the initial Makefile
go run ./.sage

# Generate targets and workflows
make update-sage-ci
```

This generates `.sage/targets.gen.go` with individual target functions based on
your configuration, giving you Makefile targets like `make go-lint`,
`make python-test`, etc.

### Run targets

```bash
# Run the default target (All)
make

# Run individual targets
make go-lint
make go-test
make python-format
```

> [!TIP]
>
> Install Makefile shell completions to see all targets in your terminal by
> typing out `make` followed by a space and then tab.

## Adding custom targets

Add a function to `.sage/sagefile.go` or a new `.go` file in `.sage/`:

```go
func MyTarget(ctx context.Context) error {
    return sg.Command(ctx, "echo", "hello").Run()
}
```

Run `go run ./.sage` to regenerate the Makefile, then use `make my-target`.

## Updating sage-ci

Either wait until the `sage-ci-sync.yml` workflow runs, or run manually:

```sh
make update-sage-ci
```

## GitHub Actions permissions

The generated workflows (e.g., `sage-ci-release.yml`, `sage-ci-sync.yml`)
require write access to create branches and pull requests. To enable this:

1. Go to your repository **Settings** → **Actions** → **General**
2. Under **Workflow permissions**, select **Read and write permissions**
3. Check **Allow GitHub Actions to create and approve pull requests**

Without these permissions, workflows will fail with a 403 "Resource not
accessible by integration" error.

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
