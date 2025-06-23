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
				middleware.ProvideAs(global_middleware.ProvideGZip),
				middleware.ProvideAs(global_middleware.ProvideCache),
				middleware.ProvideAs(global_middleware.ProvideTraceID),
				middleware.ProvideAs(global_middleware.ProvideIncomingLog),
			),
		),
	)
)
