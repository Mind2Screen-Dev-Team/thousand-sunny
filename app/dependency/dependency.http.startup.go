package dependency

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/gen/repo"
	http_middleware "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware"
	http_router "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/router"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
)

type InvokeHttpServerParam struct {
	fx.In

	Echo       *echo.Echo
	Query      *repo.Queries
	Middleware []http_middleware.Middleware `group:"global:http:middleware"`
	Router     []http_router.HttpRouter     `group:"global:http:router"`
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

	p.Echo.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, ".")
	})

	p.Echo.POST("/test", func(c echo.Context) error {
		err := c.Request().ParseMultipartForm(100 << 20)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, ".")
	})

	p.Echo.GET("/test/db", func(c echo.Context) error {
		var (
			ctx = c.Request().Context()
		)

		item, err := p.Query.FindByID(xlog.QueryName(ctx, "query.items.by.id"), uuid.MustParse("01952e0e-2283-7f7c-af3c-f4cf328aa3d6"))
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}

		return c.JSON(http.StatusOK, item)
	})

	p.Echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, http.StatusText(http.StatusOK))
	})
}
