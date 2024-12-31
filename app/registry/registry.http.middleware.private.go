package registry

import (
	private_middleware "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware/private"
	"go.uber.org/fx"
)

var (
	HttpPrivateMiddleware = fx.Options(
		fx.Module("http:server:private:middleware",
			fx.Provide(private_middleware.NewAuthJWT),
		),
	)
)