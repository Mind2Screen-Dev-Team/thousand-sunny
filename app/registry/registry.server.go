package registry

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
)

var (
	CoreServer = fx.Options(
		fx.Module("config:core:server",
			fx.Provide(dependency.ProvideHTTPServerName),
		),
	)
)
