package registry

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"go.uber.org/fx"
)

var (
	DependencyDatabase = fx.Options(
		fx.Module("dependency:database",
			fx.Provide(dependency.ProvidePostgres),
		),
	)

	DependencyDatabaseStartUp = fx.Options(
		fx.Module("dependency:database:startup",
			fx.Invoke(dependency.InvokePostgres),
		),
	)
)
