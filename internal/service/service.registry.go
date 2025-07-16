package service

import (
	"go.uber.org/fx"

	user_impl_svc "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service/user/impl"
)

var (
	Modules = fx.Module("http:services:modules",
		fx.Provide(user_impl_svc.NewExampleUserService),
	)
)
