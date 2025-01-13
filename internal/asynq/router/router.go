package asynq_router

import (
	"context"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xasynq"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type AsynqRouter interface {
	Route() xasynq.AsynqRoute
	Serve(ctx context.Context, task *asynq.Task) error
}

func RegisterAsynqAs(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(AsynqRouter)),
		fx.ResultTags(`group:"global:asynq:router"`),
	)
}
