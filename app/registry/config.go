package registry

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"go.uber.org/fx"
)

var (
	DependencyConfig = fx.Options(
		fx.Module("dependency:config",
			fx.Provide(config.ProvideConfig),
		),
	)
)
