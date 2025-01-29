package asynq_worker

import (
	"context"
	"encoding/json"

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
		m = make(map[string]any)
	)

	err := json.Unmarshal(d, &m)
	if err != nil {
		return err
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Change this as you want

	w.Write(b)
	return nil
}
