package main

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/registry"
)

func main() {
	fx.New(
		registry.Fx,
		registry.DependencyConfig,
		registry.DependencyLogger,

		// Cache
		registry.DependencyCache,
		registry.DependencyCacheStartUp,

		// Database
		registry.DependencyDatabase,
		registry.DependencyDatabaseStartUp,

		// HTTP
		registry.Http,
		registry.HttpStartUp,
		registry.HttpGlobalMiddleware,
	).Run()

	defer dependency.RotateLog()
}
