package targets

import (
	"context"

	"github.com/fredrikaverpil/sage-ci/config"
	"github.com/fredrikaverpil/sage-ci/tools/sgstylua"
	"go.einride.tech/sage/sg"
)

func luaFormat(ctx context.Context, cfg config.Config, skip SkipTargets) error {
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
