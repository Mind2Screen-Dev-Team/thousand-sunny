package module

import (
	"go.uber.org/fx"

	repo_impl "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repo/impl"
	service_impl "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service/impl"
)

var (
	UserDependencyModule = fx.Module("module:dependency:user",
		fx.Provide(
			repo_impl.NewUserCURDRepo,
			service_impl.NewUserCURDService,
		),
	)
)
