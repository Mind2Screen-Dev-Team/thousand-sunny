package user

import (
	"go.uber.org/fx"
)

var (
	RepoModules = fx.Module("repository:module:user",
		fx.Provide(NewRepo),
	)

	ServiceModules = fx.Module("service:module:user",
		fx.Provide(NewService),
	)

	HandlerModules = fx.Module("http:handler:module:user",
		fx.Provide(NewCreateHandlerFx),
		fx.Provide(NewReadAllHandlerFx),
		fx.Provide(NewReadHandlerFx),
		fx.Provide(NewUpdateHandlerFx),
		fx.Provide(NewDeleteHandlerFx),
	)
)
