package registry

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"go.uber.org/fx"
)

var (
	GlobalConfig = fx.Options(
		fx.Module("dependency:global:config",
			fx.Provide(config.ProvideConfig),
		),
	)
)
