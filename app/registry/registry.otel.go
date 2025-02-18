package registry

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"go.uber.org/fx"
)

var (
	OtelSetup = fx.Options(
		fx.Module("otel:setup",
			fx.Provide(dependency.ProvideOtelConfig),
			fx.Provide(dependency.ProvideOtelGrpcClient),
			fx.Provide(dependency.ProvideOtelResource),
			fx.Provide(dependency.ProvideOtelTracer),
			fx.Provide(dependency.ProvideOtelMetric),
		),
	)
)
