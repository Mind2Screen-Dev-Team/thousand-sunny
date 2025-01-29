package asynq_worker

import (
	"context"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xasynq"
	"github.com/hibiken/asynq"
)

type AsynqWorkerNotifyIncomingLog struct{}

func NewAsynqWorkerNotifyIncomingLog() AsynqWorkerNotifyIncomingLog {
	return AsynqWorkerNotifyIncomingLog{}
}

func (s AsynqWorkerNotifyIncomingLog) Route() xasynq.AsynqRoute {
	return xasynq.NewRoute("notify:incoming:log", 10)
}

func (s AsynqWorkerNotifyIncomingLog) Serve(ctx context.Context, task *asynq.Task) error {
	var (
		w = task.ResultWriter()
		d = task.Payload()
	)

	// Change this as you want

	w.Write(d)
	return nil
}
