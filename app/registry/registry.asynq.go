package registry

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
)

var (
	Asynq = fx.Options(
		fx.Module("asynq:server",
			fx.Provide(dependency.ProvideAsynqServerConfig),
			fx.Provide(dependency.ProvideAsynqmonOption),
			fx.Provide(dependency.ProvideAsynqMonitoringServer),
		),
	)

	AsynqPackage = fx.Options(
		fx.Module("asynq:package",
			fx.Provide(dependency.ProvideAsynqRedisConnOption),
			fx.Provide(dependency.ProvideXAsynq),
		),
	)

	AsynqStartUp = fx.Options(
		fx.Module("asynq:server:startup",
			fx.Invoke(dependency.InvokeAsynqWorkerServer),
			fx.Invoke(dependency.InvokeAsynqSchedulerServer),
			fx.Invoke(dependency.InvokeAsynqMonitoringServer),
		),
	)
)
