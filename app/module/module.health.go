package module

import (
	handler_health "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/handler/health"
	http_router "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/router"
	"go.uber.org/fx"
)

var (
	HealthHttpModule = fx.Module("module:http:health",
		fx.Provide(
			http_router.RegisterHttpAs(handler_health.NewHealthHandler),
		),
	)
)
