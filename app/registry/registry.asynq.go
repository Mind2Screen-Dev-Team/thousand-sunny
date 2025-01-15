package registry

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
)

var (
	Asynq = fx.Options(
		fx.Module("asynq:server",
			fx.Provide(dependency.ProvideAsynqRedisConnOption),
			fx.Provide(dependency.ProvideAsynqmonOption),
			fx.Provide(dependency.ProvideXAsynq),
			fx.Provide(dependency.ProvideAsynqMonitoringServer),
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
