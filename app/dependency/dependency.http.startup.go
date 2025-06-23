package dependency

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	http_middleware "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware"
	http_router "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/router"
)

type InvokeHttpServerParam struct {
	fx.In

	Cfg        config.Cfg
	Server     config.Server
	Tracer     trace.Tracer
	Echo       *echo.Echo
	Middleware []http_middleware.Middleware `group:"global:http:middleware"`
	Router     []http_router.HttpRouter     `group:"global:http:router"`
}

func InvokeHttpServer(p InvokeHttpServerParam) {
	p.Echo.Use(otelecho.Middleware(fmt.Sprintf("%s/%s", p.Cfg.App.Project, p.Server.Name)))

	sort.Slice(p.Middleware, func(i, j int) bool {
		return p.Middleware[i].Order() < p.Middleware[j].Order()
	})

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

	p.Echo.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, ".")
	})

	p.Echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, http.StatusText(http.StatusOK))
	})
}
