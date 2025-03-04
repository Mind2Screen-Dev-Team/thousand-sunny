package dependency

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xasynq"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xecho"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xresp"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xtracer"
)

func ProvideAsynqServerConfig(c config.Cfg) config.Server {
	return c.Server["asynq"]
}

func ProvideAsynqRedisConnOption(c config.Cfg) *asynq.RedisClientOpt {
	var (
		cfg   = c.Cache["redis"]
		db, _ = strconv.Atoi(cfg.DBName)
		addr  = fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
		cred  = cfg.Credential
	)

	return &asynq.RedisClientOpt{
		Addr:     addr,
		Username: cred.Username,
		Password: cred.Password,
		DB:       db,
	}
}

func ProvideAsynqmonOption(c config.Cfg, o *asynq.RedisClientOpt) *asynqmon.Options {
	var (
		cfg      = c.Server["asynq"]
		rpath, _ = cfg.Additional["asynq.route.monitoring"]
	)

	return &asynqmon.Options{
		RootPath:     rpath,
		RedisConnOpt: *o,
	}
}

func ProvideXAsynq(c config.Cfg, opt *asynq.RedisClientOpt, log *xlog.DebugLogger, loc *time.Location) *xasynq.Asynq {
	var (
		cfg, _ = c.Server["asynq"]
		all, _ = cfg.Additional["asynq.log.level"]
	)

	var (
		_ll, _ = strconv.Atoi(all)
		loglvl = asynq.LogLevel(_ll)
		logger = xasynq.NewAsynqZeroLogger(log.Logger)
	)

	return xasynq.NewAsynq(*opt, logger, loglvl, loc)
}

func ProvideAsynqMonitoringServer(c config.Cfg, l *xlog.DebugLogger, tracer trace.Tracer, lc fx.Lifecycle) *echo.Echo {
	var (
		cfg    = c.Server["asynq"]
		logger = xlog.NewLogger(l.Logger)
		srv    = echo.New()
	)

	srv.HideBanner = true
	srv.HidePort = true
	srv.Binder = new(xecho.HttpInBinder)
	srv.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		ctx, span := xtracer.Start(tracer, context.Background(), "asynq.endpoint.error.handler")
		defer span.End()

		var (
			code   = http.StatusInternalServerError
			resp   = xresp.NewRestResponse[any, any](c)
			he, ok = err.(*echo.HTTPError)
		)

		if ok {
			code = he.Code
			if err = he.Internal; err != nil {
				logger.Error(ctx, "catch internal error from handler", "error", err)
			}
		} else if err != nil {
			logger.Error(ctx, "catch error from handler", "error", err)
		}

		resp.
			StatusCode(code).
			Error(err).
			Msg(http.StatusText(code)).
			JSON()
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info(ctx, "asynq monitoring server started", "address", cfg.Address)
				if err := srv.Start(cfg.Address); err != nil && err != http.ErrServerClosed {
					logger.Error(ctx, "failed to start asynq monitoring server", "error", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			defer logger.Info(ctx, "asynq monitoring server stopped", "address", cfg.Address)
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
