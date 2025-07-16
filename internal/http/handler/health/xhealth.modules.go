package health

import (
	"go.uber.org/fx"
)

var (
	HealthHandlerModules = fx.Options(
		HealthHandlerModuleFx,
	)
)
