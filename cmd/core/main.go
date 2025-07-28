package main

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/injector"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/infra/http/middleware"
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
		injector.DatabaseStartUp,

		// Repo SQLC Generator
		injector.RepoGenerationSqlc,

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
