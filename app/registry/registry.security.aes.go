package registry

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"go.uber.org/fx"
)

var (
	SecurityEncryptionAES = fx.Options(
		fx.Module("dependency:security:aes",
			fx.Provide(dependency.ProvideEncryptionAES),
		),
	)
)
