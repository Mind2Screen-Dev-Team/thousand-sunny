package asynq_scheduler

import (
	"context"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xasynq"
	"github.com/hibiken/asynq"
)

type AsynqSchedulerExample struct{}

func NewAsynqSchedulerExample() AsynqSchedulerExample {
	return AsynqSchedulerExample{}
}

func (s AsynqSchedulerExample) Route() xasynq.AsynqRoute {
	return xasynq.NewRoute("example", 10)
}

func (s AsynqSchedulerExample) Serve(ctx context.Context, task *asynq.Task) error {
	var (
		rw = task.ResultWriter()
	)

	rw.Write([]byte(`{"status":true,"msg":"this is handler of task scheduler example"}`))
	return nil
}
