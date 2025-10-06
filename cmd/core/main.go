package main

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/injector"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/infra/http/middleware"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/infra/sdk"
)

func main() {
	fx.New(
		// Main
		injector.Fx,
		injector.GlobalConfig,
		injector.GlobalLogger,
		injector.GlobalEmail,
		injector.OtelSetup,

		// Cache
		injector.Cache,
		injector.CacheStartUp,

		// Database
		injector.Database,
		injector.GormDatabase,
		injector.GormGenDatabase,
		injector.DatabaseStartUp,

		// Repo SQLC Generator
		injector.RepoGenerationSqlc,

		// SDK
		sdk.Modules,

		// HTTP Middleware
		middleware.GlobalModules,
		middleware.PrivateModules,

		// HTTP
		injector.CoreServer,
		injector.Http,
		injector.HttpStartUp,

		// Internal Modules
		internal.RepoModules,
		internal.ServiceModules,
		internal.HandlerModules,
	).Run()

	defer dependency.RotateLog()
}
