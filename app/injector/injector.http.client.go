package injector

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
)

var (
	Resty = fx.Options(
		fx.Module("dependency:resty",
			fx.Provide(dependency.ProvideHttpClient),
		),
	)
)
