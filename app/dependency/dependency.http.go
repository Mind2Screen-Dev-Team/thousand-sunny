package dependency

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xecho"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xresp"
)

func ProvideHTTPServer(c config.Cfg, logger *xlog.DebugLogger, lc fx.Lifecycle) *echo.Echo {
	var (
		cfg = c.Server["http"]
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
				logger.Logger.Info().Str("address", cfg.Address).Msg("http server started")
				if err := srv.Start(cfg.Address); err != nil && err != http.ErrServerClosed {
					logger.Logger.Error().Err(err).Msg("failed to start http server")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			defer logger.Logger.Info().Str("address", cfg.Address).Msg("http server stopped")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
