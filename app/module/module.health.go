package module

import (
	handler_health "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/handler/health"
	http_router "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/router"
	"go.uber.org/fx"
)

var (
	HealthModule = fx.Module("module:health",
		fx.Provide(
			http_router.ProvideAs(handler_health.NewHealthHandler),
		),
	)
)
