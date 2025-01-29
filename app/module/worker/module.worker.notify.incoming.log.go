package module_worker

import (
	"go.uber.org/fx"

	asynq_router "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/asynq/router"
	asynq_worker "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/asynq/worker"
)

var (
	AsynqWorkerNotifyIncomingLogModule = fx.Module("module:asynq:worker:example",
		fx.Provide(
			asynq_router.RegisterAsynqWorkerAs(asynq_worker.NewAsynqWorkerNotifyIncomingLog),
		),
	)
)
