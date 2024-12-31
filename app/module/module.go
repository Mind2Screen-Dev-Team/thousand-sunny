package module

import "go.uber.org/fx"

func ProvideModules() fx.Option {
	return fx.Options(
		HealthModule,
		UserModule,
	)
}
