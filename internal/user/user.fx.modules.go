package user

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhuma"
)

var (
	RepoModules = fx.Module("repository:module:user",
		fx.Provide(NewRepo),
	)

	ServiceModules = fx.Module("service:module:user",
		fx.Provide(NewService),
	)

	HandlerModules = fx.Module("http:handler:module:user",
		fx.Provide(xhuma.AnnotateHandlerAs(NewCreateHandlerFx)),
		fx.Provide(xhuma.AnnotateHandlerAs(NewReadAllHandlerFx)),
		fx.Provide(xhuma.AnnotateHandlerAs(NewReadHandlerFx)),
		fx.Provide(xhuma.AnnotateHandlerAs(NewUpdateHandlerFx)),
		fx.Provide(xhuma.AnnotateHandlerAs(NewDeleteHandlerFx)),
	)
)
