package bootstrap

import (
	http_middleware "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware"

	"go.uber.org/fx"
)

func (Registry) GlobalHTTPMiddlewareProvider() Provider {
	return Provider{
		fx.Module("http-global-middleware",
			fx.Provide(
				http_middleware.As(http_middleware.ProvideCORS),
				http_middleware.As(http_middleware.ProvideRequestID),
				http_middleware.As(http_middleware.ProvideIncomingLog),
			),
		),
	}
}
