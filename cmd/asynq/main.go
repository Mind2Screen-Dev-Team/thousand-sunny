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

		// Cache
		registry.Cache,
		registry.CacheStartUp,

		// Database
		registry.Database,
		registry.DatabaseStartUp,

		// Repo SQLC Generator
		registry.RepoGenerationSqlc,

		// Modules
		module.ProvideModules,
	).Run()

	defer dependency.RotateLog()
}