package injector

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"go.uber.org/fx"
)

var (
	Cache = fx.Options(
		fx.Module("dependency:cache",
			fx.Provide(dependency.ProvideRedis),
			fx.Provide(dependency.ProvideRedisLock),
		),
	)

	CacheStartUp = fx.Options(
		fx.Module("dependency:cache:startup",
			fx.Invoke(dependency.InvokeRedis),
		),
	)
)
