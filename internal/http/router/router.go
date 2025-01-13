package http_router

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type HttpRouter interface {
	Route() (method, path string)
	Middleware() []echo.MiddlewareFunc
	Serve(c echo.Context) error
}

func RegisterHttpAs(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(HttpRouter)),
		fx.ResultTags(`group:"global:http:router"`),
	)
}
