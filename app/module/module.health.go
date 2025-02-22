package module

import (
	"go.uber.org/fx"

	handler_health "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/handler/health"
	http_router "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/router"
)

var (
	HealthHttpModule = fx.Module("module:http:health",
		fx.Provide(http_router.RegisterHttpAs(handler_health.NewHealthHandler)),
	)
)
