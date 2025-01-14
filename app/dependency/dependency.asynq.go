package dependency

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xasynq"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xecho"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xresp"
)

func ProvideAsynqServer(c config.Cfg, l *xlog.DebugLogger, lc fx.Lifecycle) *echo.Echo {
	var (
		cfg    = c.Server["asynq"]
		srv    = echo.New()
		logger = xlog.NewLogger(l.Logger)
	)

	srv.HideBanner = true
	srv.HidePort = true
	srv.Binder = new(xecho.HttpInBinder)
	srv.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var (
			code = http.StatusInternalServerError
			resp = xresp.NewRestResponse[any, any](c)
		)

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		logger.Error("catch error from handler", "error", err)
		resp.
			StatusCode(code).
			Error(err).
			Msg(http.StatusText(code)).
			JSON()
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("asynq monitoring server started", "address", cfg.Address)
				if err := srv.Start(cfg.Address); err != nil && err != http.ErrServerClosed {
					logger.Error("failed to start asynq monitoring server", "error", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			defer logger.Info("asynq monitoring server stopped", "address", cfg.Address)
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

func ProvideXAsynq(c config.Cfg, l *xlog.DebugLogger, location *time.Location) *xasynq.Asynq {
	var (
		asynqCfg = c.Server["asynq"]
		cfg      = c.Cache["redis"]
		all, _   = asynqCfg.Additional["asynq.log.level"]

		logLevel, _ = strconv.Atoi(all)
		DB, _       = strconv.Atoi(cfg.DBName)
		addr        = fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)

		opt = asynq.RedisClientOpt{
			Addr:     addr,
			Username: cfg.Credential.Username,
			Password: cfg.Credential.Password,
			DB:       DB,
		}
	)

	return xasynq.NewAsynq(
		opt,
		xasynq.NewAsynqZeroLogger(l.Logger),
		asynq.LogLevel(logLevel),
		location,
	)
}
