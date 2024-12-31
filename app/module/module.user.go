package module

import (
	"go.uber.org/fx"

	handler_user "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/handler/user"
	http_router "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/router"

	repo_impl "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repo/impl"
	service_impl "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service/impl"
)

var (
	UserModule = fx.Module("module:user",
		// service
		fx.Provide(
			repo_impl.NewUserCURDRepo,
			service_impl.NewUserCURDService,
		),

		// handler
		fx.Provide(http_router.ProvideAs(handler_user.NewUserGetDetailHandlerFx)),
		fx.Provide(http_router.ProvideAs(handler_user.NewUserCreateHandlerFx)),
	)
)
