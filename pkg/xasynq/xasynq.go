package xasynq

import (
	"context"
	"strings"
	"time"

	"github.com/hibiken/asynq"
)

type (
	AsynqCtx       string
	AsynqRouteKind string
)

func (s AsynqCtx) Str() string {
	return string(s)
}

const (
	ASYNQ_ENV AsynqCtx = "asynq:env:ctx"
)

func (s AsynqRouteKind) Is(kind AsynqRouteKind) bool {
	return s.String() == kind.String()
}

func (s AsynqRouteKind) String() string {
	return string(s)
}

const (
	ASYNQ_ROUTE_KIND_SCHEDULER AsynqRouteKind = "scheduler"
	ASYNQ_ROUTE_KIND_WORKER    AsynqRouteKind = "worker"
)

type Asynq struct {
	Client    *asynq.Client
	Inspector *asynq.Inspector
	Scheduler *asynq.Scheduler
}

func NewAsynq(opt asynq.RedisClientOpt, loc *time.Location) *Asynq {
	return &Asynq{
		Client:    asynq.NewClient(opt),
		Inspector: asynq.NewInspector(opt),
		Scheduler: asynq.NewScheduler(opt, &asynq.SchedulerOpts{Location: loc}),
	}
}

func BuildSchedulerRouteName(ctx context.Context, name string) string {
	env, _ := ctx.Value(ASYNQ_ENV).(string)
	return strings.Join([]string{env, ASYNQ_ROUTE_KIND_SCHEDULER.String(), name}, ":")
}

func BuildWorkerRouteName(ctx context.Context, name string) string {
	env, _ := ctx.Value(ASYNQ_ENV).(string)
	return strings.Join([]string{env, ASYNQ_ROUTE_KIND_WORKER.String(), name}, ":")
}

type (
	// Router name is already include env and type of route. ex:
	//	- development:scheduler:<your_route_name>
	//	- development:worker:<your_route_name>
	AsynqRoute struct {
		Name   string
		Kind   AsynqRouteKind
		Worker int
	}
)

func NewSchedulerRoute(name string, worker int) AsynqRoute {
	return AsynqRoute{Name: name, Kind: ASYNQ_ROUTE_KIND_SCHEDULER, Worker: worker}
}

func NewWorkerRoute(name string, worker int) AsynqRoute {
	return AsynqRoute{Name: name, Kind: ASYNQ_ROUTE_KIND_WORKER, Worker: worker}
}
