package registry

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
)

var (
	Http = fx.Options(
		fx.Module("http:server",
			fx.Provide(dependency.ProvideHTTPServer),
		),
	)

	HttpStartUp = fx.Options(
		fx.Module("http:server:startup",
			fx.Invoke(dependency.InvokeHttpServer),
		),
	)
)
