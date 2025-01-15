package asynq_router

import (
	"context"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xasynq"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type AsynqWorkerRouter interface {
	Route() xasynq.AsynqRoute
	Serve(ctx context.Context, task *asynq.Task) error
}

func RegisterAsynqWorkerAs(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(AsynqWorkerRouter)),
		fx.ResultTags(`group:"global:asynq:worker:router"`),
	)
}

type AsynqSchedulerRouter interface {
	Route() xasynq.AsynqRoute
	Serve(ctx context.Context, task *asynq.Task) error
}

func RegisterAsynqSchedulerAs(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(AsynqSchedulerRouter)),
		fx.ResultTags(`group:"global:asynq:scheduler:router"`),
	)
}
