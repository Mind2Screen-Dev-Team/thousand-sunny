package xasynq

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
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

func NewAsynq(opt asynq.RedisClientOpt, logger asynq.Logger, logLvl asynq.LogLevel, loc *time.Location) *Asynq {
	return &Asynq{
		Client:    asynq.NewClient(opt),
		Inspector: asynq.NewInspector(opt),
		Scheduler: asynq.NewScheduler(opt, &asynq.SchedulerOpts{Location: loc, Logger: logger, LogLevel: logLvl}),
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

type AsynqZeroLogger struct {
	logger zerolog.Logger
}

func NewAsynqZeroLogger(logger zerolog.Logger) asynq.Logger {
	return &AsynqZeroLogger{logger}
}

// Debug implements asynq.Logger.
func (a *AsynqZeroLogger) Debug(args ...any) {
	msg, ok := args[0].(string)
	if !ok {
		msg = fmt.Sprintf("%v", args[0])
	}

	a.logger.Debug().Fields(args[1:]).Msg(msg)
}

// Error implements asynq.Logger.
func (a *AsynqZeroLogger) Error(args ...any) {
	msg, ok := args[0].(string)
	if !ok {
		msg = fmt.Sprintf("%v", args[0])
	}

	a.logger.Error().Fields(args[1:]).Msg(msg)
}

// Fatal implements asynq.Logger.
func (a *AsynqZeroLogger) Fatal(args ...any) {
	msg, ok := args[0].(string)
	if !ok {
		msg = fmt.Sprintf("%v", args[0])
	}

	a.logger.Fatal().Fields(args[1:]).Msg(msg)
}

// Info implements asynq.Logger.
func (a *AsynqZeroLogger) Info(args ...any) {
	msg, ok := args[0].(string)
	if !ok {
		msg = fmt.Sprintf("%v", args[0])
	}

	a.logger.Info().Fields(args[1:]).Msg(msg)
}

// Warn implements asynq.Logger.
func (a *AsynqZeroLogger) Warn(args ...any) {
	msg, ok := args[0].(string)
	if !ok {
		msg = fmt.Sprintf("%v", args[0])
	}

	a.logger.Warn().Fields(args[1:]).Msg(msg)
}
