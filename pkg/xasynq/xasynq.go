package xasynq

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

// # Utilitty

const (
	ASYNQ_ENV        AsynqCtx = "asynq:env:ctx"
	ASYNQ_ROUTE_KIND AsynqCtx = "asynq:route:kind:ctx"
)

const (
	ASYNQ_ROUTE_KIND_SCHEDULER AsynqRouteKind = "scheduler"
	ASYNQ_ROUTE_KIND_WORKER    AsynqRouteKind = "worker"
)

type (
	AsynqCtx       string
	AsynqRouteKind string
)

func (s AsynqCtx) Str() string {
	return string(s)
}

func (s AsynqRouteKind) Is(kind AsynqRouteKind) bool {
	return s.String() == kind.String()
}

func (s AsynqRouteKind) String() string {
	return string(s)
}

// # Asynq Collection Client

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

// # Builder

type (
	// Router name is already include env and type of route. ex:
	//	- development:scheduler:<your_route_name>
	//	- development:worker:<your_route_name>
	AsynqRoute struct {
		Name        string
		Concurrency int
	}
)

// Router name is already include env and type of route. ex:
//   - development:scheduler:<your_route_name>
//   - development:worker:<your_route_name>
func NewRoute(name string, concurrency int) AsynqRoute {
	return AsynqRoute{Name: name, Concurrency: concurrency}
}

// This will help you build route name:
//
// Note:
//   - The context design for get a value
//   - The env and kind of route (worker / scheduler) is set on base context
//
// Example:
//   - development:worker:<name>
func BuildRouteName(ctx context.Context, name string) string {
	var (
		env, _  = ctx.Value(ASYNQ_ENV).(string)
		kind, _ = ctx.Value(ASYNQ_ROUTE_KIND).(string)
	)
	return strings.Join([]string{env, kind, name}, ":")
}

// This will help you build route name, ex:
//   - development:worker:<names>...
func BuildWorkerRouteName(env string, names ...string) string {
	var (
		kind  = ASYNQ_ROUTE_KIND_WORKER.String()
		slice = []string{env, kind}
	)
	return strings.Join(append(slice, names...), ":")
}

// This will help you build route name, ex:
//   - development:scheduler:<names>...
func BuildSchedulerRouteName(env string, names ...string) string {
	var (
		kind  = ASYNQ_ROUTE_KIND_SCHEDULER.String()
		slice = []string{env, kind}
	)
	return strings.Join(append(slice, names...), ":")
}

// # Logger

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
