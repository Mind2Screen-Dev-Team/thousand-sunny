package handler

import (
	"go.uber.org/fx"

	handler_health "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/handler/health"
)

var (
	Modules = fx.Module("http:handler:modules",
		handler_health.HealthHandlerModuleFx,
		// add more handler here
	)
)
