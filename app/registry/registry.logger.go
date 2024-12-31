package registry

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"go.uber.org/fx"
)

var (
	GlobalLogger = fx.Options(
		fx.Module("dependency:global:logger",
			fx.Provide(dependency.ProvideDebugLogger),
			fx.Provide(dependency.ProvideIoLogger),
			fx.Provide(dependency.ProvideTrxLogger),
		),
	)
)
