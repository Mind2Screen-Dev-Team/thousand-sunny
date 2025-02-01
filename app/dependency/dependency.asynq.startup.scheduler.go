package dependency

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xasynq"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"

	asynq_router "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/asynq/router"
)

type InvokeAsynqSchedulerServerParam struct {
	fx.In

	Cfg          config.Cfg
	Lfc          fx.Lifecycle
	Log          *xlog.DebugLogger
	Asynq        *xasynq.Asynq
	RedisConnOpt *asynq.RedisClientOpt
	Router       []asynq_router.AsynqSchedulerRouter `group:"global:asynq:scheduler:router"`
}

func InvokeAsynqSchedulerServer(p InvokeAsynqSchedulerServerParam) error {
	if p.Log == nil {
		return errors.New("field 'Log' with type '*xlog.DebugLogger' is not provided")
	}

	if p.Asynq == nil {
		return errors.New("field 'Asynq' with type '*xasynq.Asynq' is not provided")
	}

	if p.RedisConnOpt == nil {
		return errors.New("field 'RedisConnOpt' with type '*asynq.RedisClientOpt' is not provided")
	}

	var (
		env     = p.Cfg.App.Env
		acfg, _ = p.Cfg.Server["asynq"]
		all, _  = acfg.Additional["asynq.log.level"]

		ll, _  = strconv.Atoi(all)
		logger = xlog.NewLogger(p.Log.Logger)
		kind   = xasynq.ASYNQ_ROUTE_KIND_SCHEDULER.String()
	)

	var (
		mux = asynq.NewServeMux()
		cfg = asynq.Config{
			Queues:   make(map[string]int),
			Logger:   xasynq.NewAsynqZeroLogger(p.Log.Logger),
			LogLevel: asynq.LogLevel(ll),
			BaseContext: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, xasynq.ASYNQ_ENV, env)
				ctx = context.WithValue(ctx, xasynq.ASYNQ_ROUTE_KIND, kind)
				return ctx
			},
		}
	)

	for _, router := range p.Router {
		var (
			route = router.Route()
			slice = []string{env, kind, route.Name}
			queue = strings.ToLower(strings.Join(slice, ":"))
		)

		if _, exist := cfg.Queues[queue]; !exist {
			cfg.Concurrency += route.Concurrency
			cfg.Queues[queue] = route.Concurrency
			mux.HandleFunc(queue, router.Serve)
		}
	}

	var (
		server = asynq.NewServer(*p.RedisConnOpt, cfg)
	)

	p.Lfc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			go func() {
				logger.Info("asynq scheduler started")
				if err := server.Start(mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("failed to start asynq scheduler server", "error", err)
				}
			}()

			go func() {
				logger.Info("asynq cron scheduler started")
				if err := p.Asynq.Scheduler.Run(); err != nil {
					logger.Error("failed to start asynq cron scheduler", "error", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {

			server.Shutdown()
			logger.Info("asynq scheduler server stopped")

			p.Asynq.Scheduler.Shutdown()
			logger.Info("asynq cron scheduler stopped")

			return nil
		},
	})

	logger.Info("asynq information of concurrency and queue", "kind", "scheduler", "concurrency", cfg.Concurrency, "queue", cfg.Queues)
	return nil
}
