package dependency

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xasynq"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"

	asynq_router "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/asynq/router"
)

type InvokeAsynqServerParam struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    config.Cfg
	Echo      *echo.Echo
	logger    *xlog.DebugLogger
	XAsynq    *xasynq.Asynq
	Router    []asynq_router.AsynqRouter `group:"global:asynq:router"`
}

func InvokeAsynqServer(p InvokeAsynqServerParam) {
	var (
		env   = p.Config.App.Env
		cfg   = p.Config.Cache["redis"]
		DB, _ = strconv.Atoi(cfg.DBName)
		ctx   = context.WithValue(context.Background(), xasynq.ASYNQ_ENV, env)
	)

	var (
		redisOpt = asynq.RedisClientOpt{
			Addr:     cfg.Address,
			Username: cfg.Credential.Username,
			Password: cfg.Credential.Password,
			DB:       DB,
		}

		asynqTaskMux     = asynq.NewServeMux()
		asynqScheduleMux = asynq.NewServeMux()

		asynqTaskCfg = asynq.Config{
			Queues: make(map[string]int),
			BaseContext: func() context.Context {
				return ctx
			},
		}

		asynqScheduleCfg = asynq.Config{
			Queues: make(map[string]int),
			BaseContext: func() context.Context {
				return ctx
			},
		}
	)

	var (
		_rootpath   = "/monitoring/tasks"
		asynqmonCfg = asynqmon.Options{
			RootPath:     _rootpath,
			RedisConnOpt: redisOpt,
		}
	)

	for _, router := range p.Router {
		var (
			route = router.Route()
			slice = []string{env, route.Kind.String(), route.Name}
			queue = strings.ToLower(strings.Join(slice, ":"))
		)

		switch route.Kind {
		case xasynq.ASYNQ_ROUTE_KIND_WORKER:
			{
				if _, exist := asynqTaskCfg.Queues[queue]; !exist {
					asynqTaskCfg.Concurrency += route.Worker
					asynqTaskCfg.Queues[queue] = route.Worker
					asynqTaskMux.HandleFunc(queue, router.Serve)
				}
			}
		case xasynq.ASYNQ_ROUTE_KIND_SCHEDULER:
			{
				if _, exist := asynqScheduleCfg.Queues[queue]; !exist {
					asynqScheduleCfg.Concurrency += route.Worker
					asynqScheduleCfg.Queues[queue] = route.Worker
					asynqScheduleMux.HandleFunc(queue, router.Serve)
				}
			}
		default:
			continue
		}
	}

	var (
		_asynqmon            = asynqmon.New(asynqmonCfg)
		_asynqTaskServer     = asynq.NewServer(redisOpt, asynqTaskCfg)
		_asynqScheduleServer = asynq.NewServer(redisOpt, asynqScheduleCfg)
		_handler             = func(c echo.Context) error {
			return c.String(http.StatusOK, ".")
		}
	)

	// add lifecycle for asynq server
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			go func() {
				p.logger.Logger.Info().Msg("asynq tasks server started")
				if err := _asynqTaskServer.Start(asynqTaskMux); err != nil && !errors.Is(err, http.ErrServerClosed) {
					p.logger.Logger.Error().Err(err).Msg("failed to start asynq tasks server")
				}
			}()

			go func() {
				p.logger.Logger.Info().Msg("asynq schedule server started")
				if err := _asynqScheduleServer.Start(asynqScheduleMux); err != nil && !errors.Is(err, http.ErrServerClosed) {
					p.logger.Logger.Error().Err(err).Msg("failed to start asynq schedule server")
				}
			}()

			go func() {
				p.logger.Logger.Info().Msg("asynq cron scheduler started")
				if err := p.XAsynq.Scheduler.Run(); err != nil {
					p.logger.Logger.Error().Err(err).Msg("failed to start asynq cron scheduler")
				}
			}()

			return nil
		}, OnStop: func(ctx context.Context) error {

			go func() {
				defer p.logger.Logger.Info().Msg("asynq tasks server stopped")
				_asynqTaskServer.Shutdown()
			}()

			go func() {
				defer p.logger.Logger.Info().Msg("asynq schedule server stopped")
				_asynqScheduleServer.Shutdown()
			}()

			go func() {
				defer p.logger.Logger.Info().Msg("asynq cron scheduler stopped")
				p.XAsynq.Scheduler.Shutdown()
			}()

			return nil
		},
	})

	// routers
	p.Echo.Use(middleware.CORS())
	p.Echo.Use(middleware.Recover())
	p.Echo.Any(fmt.Sprintf("%s/*", _rootpath), echo.WrapHandler(_asynqmon))
	p.Echo.GET("/ping", _handler)
	p.Echo.GET("/", _handler)
}
