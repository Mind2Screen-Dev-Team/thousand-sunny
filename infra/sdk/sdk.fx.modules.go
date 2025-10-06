package sdk

import (
	"go.uber.org/fx"
)

var (
	Modules = fx.Options(
		fx.Module("sdk:modules"),
	)
)
