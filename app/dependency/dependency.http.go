package dependency

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xecho"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xresp"

	httpin_integration "github.com/ggicci/httpin/integration"
)

func ProvideHTTPServerName(c config.Cfg) config.Server {
	return c.Server["http"]
}

func ProvideHTTPServer(c config.Cfg, logger *xlog.DebugLogger, lc fx.Lifecycle) *echo.Echo {
	var (
		ctx = context.Background()
		log = xlog.NewLogger(logger.Logger)
		cfg = c.Server["http"]
		srv = echo.New()
	)

	httpin_integration.UseEchoRouter("path", srv)

	if v, ok := cfg.Additional["idle.timeout"]; ok {
		if n, err := strconv.ParseInt(v, 10, 64); n >= 0 && err == nil {
			srv.Server.IdleTimeout = time.Duration(n) * time.Second
		} else if err != nil {
			log.Error(ctx, fmt.Sprintf("failed to parse http additional config: '%s'", "idle.timeout"), "err", fmt.Sprintf("%+v", err))
		}
	}

	if v, ok := cfg.Additional["write.timeout"]; ok {
		if n, err := strconv.ParseInt(v, 10, 64); n >= 0 && err == nil {
			srv.Server.WriteTimeout = time.Duration(n) * time.Second
		} else if err != nil {
			log.Error(ctx, fmt.Sprintf("failed to parse http additional config: '%s'", "write.timeout"), "err", fmt.Sprintf("%+v", err))
		}
	}

	if v, ok := cfg.Additional["read.timeout"]; ok {
		if n, err := strconv.ParseInt(v, 10, 64); n >= 0 && err == nil {
			srv.Server.ReadTimeout = time.Duration(n) * time.Second
		} else if err != nil {
			log.Error(ctx, fmt.Sprintf("failed to parse http additional config: '%s'", "read.timeout"), "err", fmt.Sprintf("%+v", err))
		}
	}

	if v, ok := cfg.Additional["read.header.timeout"]; ok {
		if n, err := strconv.ParseInt(v, 10, 64); n >= 0 && err == nil {
			srv.Server.ReadHeaderTimeout = time.Duration(n) * time.Second
		} else if err != nil {
			log.Error(ctx, fmt.Sprintf("failed to parse http additional config: '%s'", "read.header.timeout"), "err", fmt.Sprintf("%+v", err))
		}
	}

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

	srv.RouteNotFound("/*", func(c echo.Context) error {
		resp := xresp.NewRestResponse[any, any](c)
		return resp.
			StatusCode(http.StatusNotFound).
			Code(http.StatusNotFound).
			Error("route not found").
			Msg(http.StatusText(http.StatusNotFound)).
			JSON()
	})

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
