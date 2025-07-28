package internal

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/health"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/user"
)

var (
	RepoModules = fx.Options(
		user.RepoModules,
	)

	ServiceModules = fx.Options(
		user.ServiceModules,
	)

	HandlerModules = fx.Options(
		health.HandlerModules,
		user.HandlerModules,
	)
)
