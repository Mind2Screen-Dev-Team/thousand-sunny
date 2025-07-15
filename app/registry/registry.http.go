package registry

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
)

var (
	Http = fx.Options(
		fx.Module("http:server",
			fx.Provide(dependency.ProvideHumaConfig),
			fx.Provide(dependency.ProvideFiberConfig),
			fx.Provide(dependency.ProvideFiber),
			fx.Provide(dependency.ProvideHumaFiber),
			fx.Provide(dependency.ProvideHTTPServerName),
		),
	)

	HttpStartUp = fx.Options(
		fx.Module("http:server:startup",
			fx.Invoke(dependency.InvokeHTTPServer),
		),
	)
)
