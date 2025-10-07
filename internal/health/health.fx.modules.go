package health

import (
	"go.uber.org/fx"
)

var (
	HandlerModules = fx.Module("http:handler:module:health",
		fx.Provide(NewHealthHandlerFx),
	)
)
