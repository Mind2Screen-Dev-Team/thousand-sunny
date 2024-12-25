package dependency

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	http_middleware "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
)

func ProvideEchoServer(c config.Cfg, ms []http_middleware.Middleware, logger *xlog.DebugLogger, lc fx.Lifecycle) *echo.Echo {
	srv := echo.New()

	// # load global middleware
	for _, m := range ms {
		srv.Use(m.Serve)
	}

	srv.HideBanner = true
	srv.HidePort = true

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Logger.Info().Str("address", c.Http.Address).Msg("http server started")
				err := srv.Start(c.Http.Address)
				if err != nil && err != http.ErrServerClosed {
					log.Fatal(err.Error())
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			defer logger.Logger.Info().Str("address", c.Http.Address).Msg("http server stopped")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

func InitEchoMiddleware(ms []http_middleware.Middleware) []echo.MiddlewareFunc {
	fn := make([]echo.MiddlewareFunc, len(ms))
	for _, m := range ms {
		fn = append(fn, m.Serve)
	}
	return fn
}

func InitEchoServer(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
}
