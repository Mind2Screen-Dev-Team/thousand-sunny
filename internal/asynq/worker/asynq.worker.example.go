package asynq_worker

import (
	"context"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xasynq"
	"github.com/hibiken/asynq"
)

type AsynqWorkerExample struct{}

func NewAsynqWorkerExample() AsynqWorkerExample {
	return AsynqWorkerExample{}
}

func (s AsynqWorkerExample) Route() xasynq.AsynqRoute {
	return xasynq.NewWorkerRoute("example", 10)
}

func (s AsynqWorkerExample) Serve(ctx context.Context, task *asynq.Task) error {
	var (
		rw = task.ResultWriter()
	)

	rw.Write([]byte(`{"status":true,"msg":"this is handler of task worker example"}`))
	return nil
}
