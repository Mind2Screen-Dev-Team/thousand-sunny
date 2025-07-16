package main

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/registry"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/handler"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/provider"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repository"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service"
)

func main() {
	fx.New(
		// Main
		registry.Fx,
		registry.GlobalConfig,
		registry.GlobalLogger,
		registry.GlobalEmail,
		registry.OtelSetup,

		// Cache
		registry.Cache,
		registry.CacheStartUp,

		// Database
		registry.Database,
		registry.DatabaseStartUp,

		// Repo SQLC Generator
		registry.RepoGenerationSqlc,

		// HTTP Middleware
		middleware.GlobalModules,
		middleware.PrivateModules,

		// HTTP
		registry.CoreServer,
		registry.Http,
		registry.HttpStartUp,

		// Modules
		provider.Modules,
		repository.Modules,
		service.Modules,
		handler.Modules,
	).Run()

	defer dependency.RotateLog()
}
