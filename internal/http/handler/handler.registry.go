package handler

import (
	"go.uber.org/fx"

	hdr_health "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/handler/health"
	hdr_user "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/handler/user"
)

var (
	Modules = fx.Module("http:handler:modules",
		hdr_health.HealthHandlerModules,
		hdr_user.ExampleUserHandlerModules,
	)
)
