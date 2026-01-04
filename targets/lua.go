package targets

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/tools/sgstylua"
	"github.com/fredrikaverpil/sage-ci/workflows"
	"go.einride.tech/sage/sg"
)

// LuaFormat applies Lua formatting using stylua.
func LuaFormat(ctx context.Context, cfg workflows.Config, skip SkipTargets) error {
	for _, module := range cfg.LuaModules {
		if skip.ShouldSkip("LuaFormat", module) {
			continue
		}
		sg.Logger(ctx).Printf("applying stylua format in %s...", module)
		cmd := sgstylua.Command(ctx, ".")
		cmd.Dir = sg.FromGitRoot(module)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
