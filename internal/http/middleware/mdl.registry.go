package middleware

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhuma"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware/global"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware/private"
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
				xhuma.AnnotateGlobalMiddlewareAs(global.ProvideOtel),
				xhuma.AnnotateGlobalMiddlewareAs(global.ProvideTraceID),
				xhuma.AnnotateGlobalMiddlewareAs(global.ProvideHelmet),
				xhuma.AnnotateGlobalMiddlewareAs(global.ProvideCORS),
				xhuma.AnnotateGlobalMiddlewareAs(global.ProvideIncomingLog),
				xhuma.AnnotateGlobalMiddlewareAs(global.ProvideCache),
				xhuma.AnnotateGlobalMiddlewareAs(global.ProvideCompress),
				xhuma.AnnotateGlobalMiddlewareAs(global.ProvideMonitor),
				xhuma.AnnotateGlobalMiddlewareAs(global.ProvideFavicon),
			),
		),
	)
)

var (
	PrivateModules = fx.Options(
		fx.Module("http:server:private:middleware",
			fx.Provide(private.NewAuthJWT),
		),
	)
)
