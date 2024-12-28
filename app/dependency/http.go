package dependency

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	m "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
)

func ProvideHTTPServer(c config.Cfg, logger *xlog.DebugLogger, lc fx.Lifecycle) *echo.Echo {
	var (
		cfg = c.Server["http"]
		srv = echo.New()
	)

	srv.HideBanner = true
	srv.HidePort = true

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Logger.Info().Str("address", cfg.Address).Msg("http server started")
				err := srv.Start(cfg.Address)
				if err != nil && err != http.ErrServerClosed {
					log.Fatal(err.Error())
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

func InvokeHttpServer(ms []m.Middleware, e *echo.Echo) {
	for _, v := range ms {
		e.Use(v.Serve)
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
}
