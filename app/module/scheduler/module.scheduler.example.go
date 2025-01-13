package module_scheduler

import (
	"go.uber.org/fx"

	asynq_router "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/asynq/router"
	asynq_scheduler "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/asynq/scheduler"
)

var (
	AsynqSchedulerExampleModule = fx.Module("module:asynq:scheduler:example",
		fx.Provide(
			asynq_router.RegisterAsynqAs(asynq_scheduler.NewAsynqSchedulerExample),
		),
	)
)
