package injector

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
)

var (
	OtelSetup = fx.Options(
		fx.Module("otel:setup",
			fx.Provide(dependency.ProvideOtelConfig),
			fx.Provide(dependency.ProvideOtelGrpcClient),
			fx.Provide(dependency.ProvideOtelResource),
			fx.Provide(dependency.ProvideOtelTracer),
			fx.Provide(dependency.ProvideOtelMetric),
			fx.Provide(dependency.ProvideOtelLog),
		),
	)
)
