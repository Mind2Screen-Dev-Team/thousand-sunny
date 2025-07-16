package repository

import (
	"go.uber.org/fx"

	user_impl_repo "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repository/user/impl"
)

var (
	Modules = fx.Module("http:repositories:modules",
		fx.Provide(user_impl_repo.NewExampleUserRepo),
	)
)
