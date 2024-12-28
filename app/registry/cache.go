package registry

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"go.uber.org/fx"
)

var (
	DependencyCache = fx.Options(
		fx.Module("dependency:cache",
			fx.Provide(dependency.ProvideRedis),
		),
	)

	DependencyCacheStartUp = fx.Options(
		fx.Module("dependency:cache:startup",
			fx.Invoke(dependency.InvokeRedis),
		),
	)
)
