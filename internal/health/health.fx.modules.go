package health

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhuma"
)

var (
	HandlerModules = fx.Module("http:handler:module:health",
		fx.Provide(xhuma.AnnotateHandlerAs(NewHealthHandlerFx)),
	)
)
