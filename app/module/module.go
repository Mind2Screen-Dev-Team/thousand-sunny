package module

import (
	"go.uber.org/fx"

	module_scheduler "github.com/Mind2Screen-Dev-Team/thousand-sunny/app/module/scheduler"
	module_worker "github.com/Mind2Screen-Dev-Team/thousand-sunny/app/module/worker"
)

var (
	ProvideHttpModules = fx.Options(
		HealthHttpModule,
		UserHttpModule,
	)

	ProvideAsynqModules = fx.Options(
		module_scheduler.AsynqSchedulerExampleModule,
		module_worker.AsynqWorkerExampleModule,
	)
)
