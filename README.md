# sage-tools

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
