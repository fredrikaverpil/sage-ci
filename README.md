# sage-ci (powered by [Sage](https://github.com/einride/sage) 🌿)

This repo holds tools and configurations which I want synced across several
projects.

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
import (
    "github.com/fredrikaverpil/sage-ci/targets"
    "github.com/fredrikaverpil/sage-ci/config"
)

var cfg = config.Config{
    GoModules: []string{"."},
    PythonModules: []string{"tests"},
    LuaModules: []string{"lua"},
    Platform: config.PlatformGitHub,
}

var skip = targets.SkipTargets{}

func All(ctx context.Context) error {
    sg.Deps(ctx, GenerateWorkflows)
    sg.SerialDeps(ctx, RunSerial)
    sg.Deps(ctx, RunParallel)
    return targets.GitDiffCheck(ctx)
}
```

See [config/config.go](config/config.go) for all configuration options.

You can add custom targets to `sagefile.go` or create additional `.go` files in
`.sage/`. Sage-ci provides opinionated targets in `RunSerial` and `RunParallel`.

### Generate and run Makefile

```bash
# Generate the Makefile
go run ./.sage

# Generate workflows
make generate-workflows

# Run Makefile (runs the All() function of sagefile.go)
make
```

> [!TIP]
>
> Install Makefile shell completions to see all targets in your terminal by
> typing out `make` followed by a space and then tab.

## Updating sage-ci

Either wait until the `sage-ci-sync.yml` workflow runs, or run manually:

```sh
make update-sage-ci
```

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
