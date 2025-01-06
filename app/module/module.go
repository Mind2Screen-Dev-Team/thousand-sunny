package module

import "go.uber.org/fx"

var (
	ProvideModules = fx.Options(
		HealthModule,
		UserModule,
	)
)
