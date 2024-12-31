package registry

import (
	"go.uber.org/fx"

	middleware "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware"
	global_middleware "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware/global"
)

var (
	HttpGlobalMiddleware = fx.Options(
		fx.Module("http:server:global:middleware",
			fx.Provide(
				middleware.ProvideAs(global_middleware.ProvideCORS),
				middleware.ProvideAs(global_middleware.ProvideRequestID),
				middleware.ProvideAs(global_middleware.ProvideIncomingLog),
			),
		),
	)
)
