package registry

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	middleware "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware"
)

var (
	Http = fx.Options(
		fx.Module("http:server",
			fx.Provide(dependency.ProvideHTTPServer),
		),
	)

	HttpStartUp = fx.Options(
		fx.Module("http:server:startup",
			fx.Invoke(middleware.To(dependency.InvokeHttpServer)),
		),
	)

	HttpGlobalMiddleware = fx.Options(
		fx.Module("http:server:global:middleware",
			fx.Provide(
				middleware.As(middleware.ProvideCORS),
				middleware.As(middleware.ProvideRequestID),
				middleware.As(middleware.ProvideIncomingLog),
			),
		),
	)
)
