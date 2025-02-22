package module

import (
	"go.uber.org/fx"

	handler_user "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/handler/user"
	http_router "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/router"
)

var (
	UserHttpModule = fx.Module("module:http:user",
		// handler
		fx.Provide(http_router.RegisterHttpAs(handler_user.NewUserGetDetailHandlerFx)),
		fx.Provide(http_router.RegisterHttpAs(handler_user.NewUserCreateHandlerFx)),
	)
)
