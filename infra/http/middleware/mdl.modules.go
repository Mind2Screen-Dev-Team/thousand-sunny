package middleware

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhuma"
)

var (
	GlobalMiddlewareOrder = map[string]int{
		"otel.http":    0,
		"trace.id":     1,
		"helmet":       3,
		"cors":         4,
		"incoming.log": 5,
		"cache":        6,
		"compress":     7,
		"monitor":      8,
		"favicon":      9,
	}

	GlobalModules = fx.Options(
		fx.Module("http:server:global:middleware",
			fx.Provide(
				xhuma.AnnotateGlobalMiddlewareAs(ProvideOtel),
				xhuma.AnnotateGlobalMiddlewareAs(ProvideTraceID),
				xhuma.AnnotateGlobalMiddlewareAs(ProvideHelmet),
				xhuma.AnnotateGlobalMiddlewareAs(ProvideCORS),
				xhuma.AnnotateGlobalMiddlewareAs(ProvideIncomingLog),
				xhuma.AnnotateGlobalMiddlewareAs(ProvideCache),
				xhuma.AnnotateGlobalMiddlewareAs(ProvideCompress),
				xhuma.AnnotateGlobalMiddlewareAs(ProvideMonitor),
				xhuma.AnnotateGlobalMiddlewareAs(ProvideFavicon),
			),
		),
	)
)

var (
	PrivateModules = fx.Options(
		fx.Module("http:server:private:middleware",
			fx.Provide(NewPrivateAuthJWT),
		),
	)
)
