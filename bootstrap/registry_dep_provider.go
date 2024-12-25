package bootstrap

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/bootstrap/dependency"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"go.uber.org/fx"
)

func (Registry) DependencyProvider() Provider {
	return Provider{
		fx.Module("app-core",
			fx.Provide(
				config.ProvideConfig,
				dependency.ProvideDebugLogger,
				dependency.ProvideIOLogger,
				dependency.ProvideTRXLogger,
				dependency.ProvidePGxDB,
				dependency.ProvideRedisDB,
			),
			fx.Provide(
				fx.Annotate(
					dependency.ProvideEchoServer,
					fx.ParamTags("", `group:"middlewares"`),
				),
			),
		),
	}
}
