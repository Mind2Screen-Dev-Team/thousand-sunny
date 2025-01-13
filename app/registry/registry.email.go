package registry

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"go.uber.org/fx"
)

var (
	GlobalEmail = fx.Options(
		fx.Module("dependency:global:email",
			fx.Provide(dependency.ProvideXGmail),
		),
	)
)
