package main

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/module"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/registry"
)

func main() {
	fx.New(
		// Main
		registry.Fx,
		registry.GlobalConfig,
		registry.GlobalLogger,
		registry.GlobalEmail,

		// Cache
		registry.AsynqPackage,
		registry.Cache,
		registry.CacheStartUp,

		// Database
		registry.Database,
		registry.DatabaseStartUp,

		// Repo SQLC Generator
		registry.RepoGenerationSqlc,

		// HTTP
		registry.Http,
		registry.HttpStartUp,
		registry.HttpGlobalMiddleware,
		registry.HttpPrivateMiddleware,

		// Modules
		module.ProvideHttpModules,
	).Run()

	defer dependency.RotateLog()
}
