package dependency

import (
	"context"
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

func ProvideAsynqServer(c config.Cfg, logger *xlog.DebugLogger, lc fx.Lifecycle) *echo.Echo {
	var (
		cfg = c.Server["asynq"]
		srv = echo.New()
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

		logger.Logger.Error().Err(err).Msg("catch error from handler")
		resp.
			StatusCode(code).
			Error(err).
			Msg(http.StatusText(code)).
			JSON()
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Logger.Info().Str("address", cfg.Address).Msg("asynq monitoring server started")
				if err := srv.Start(cfg.Address); err != nil && err != http.ErrServerClosed {
					logger.Logger.Error().Err(err).Msg("failed to start asynq monitoring server")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			defer logger.Logger.Info().Str("address", cfg.Address).Msg("asynq monitoring server stopped")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

func ProvideXAsynq(c config.Cfg, location *time.Location) *xasynq.Asynq {
	var (
		cfg   = c.Cache["redis"]
		DB, _ = strconv.Atoi(cfg.DBName)
		opt   = asynq.RedisClientOpt{
			Addr:     cfg.Address,
			Username: cfg.Credential.Username,
			Password: cfg.Credential.Password,
			DB:       DB,
		}
	)

	return xasynq.NewAsynq(opt, location)
}
