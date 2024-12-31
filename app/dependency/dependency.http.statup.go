package dependency

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	http_middleware "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware"
	http_router "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/router"
)

type InvokeHttpServerParam struct {
	fx.In

	Echo       *echo.Echo
	Middleware []http_middleware.Middleware `group:"global:http:middleware"`
	Router     []http_router.Router         `group:"global:http:router"`
}

func InvokeHttpServer(p InvokeHttpServerParam) {
	for _, v := range p.Middleware {
		// global middlewares
		p.Echo.Use(v.Serve)
	}

	for _, r := range p.Router {
		var (
			_method, _path = r.Route()
			mdls           = r.Middleware()
		)

		// route handlers
		if len(mdls) > 0 {
			p.Echo.Add(_method, _path, r.Serve, mdls...)
			continue
		}

		p.Echo.Add(_method, _path, r.Serve)
	}

	p.Echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, ".")
	})
}
