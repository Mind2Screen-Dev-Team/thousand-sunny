package registry

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"go.uber.org/fx"
)

var (
	Dependency = fx.Options(
		fx.Module("dependency:config",
			fx.Provide(config.ProvideConfig),
		),
		fx.Module("dependency:logger",
			fx.Provide(dependency.ProvideDebugLogger),
			fx.Provide(dependency.ProvideIoLogger),
			fx.Provide(dependency.ProvideTrxLogger),
		),
		fx.Module("dependency:cache",
			fx.Provide(dependency.ProvideRedis),
		),
		fx.Module("dependency:database",
			fx.Provide(dependency.ProvidePostgres),
		),
	)

	DependencyStartUp = fx.Options(
		fx.Module("dependency:startup",
			fx.Invoke(dependency.InvokeRedis),
			fx.Invoke(dependency.InvokePostgres),
		),
	)
)
