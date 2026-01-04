package targets

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/config"
	"github.com/fredrikaverpil/sage-ci/tools/sgstylua"
	"go.einride.tech/sage/sg"
)

// LuaFormat runs stylua for all configured Lua modules.
func LuaFormat(ctx context.Context, cfg config.Config) error {
	for _, module := range cfg.LuaModules {
		if cfg.SkipTargets.ShouldSkip("LuaFormat", module) {
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
