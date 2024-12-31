package http_router

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Router interface {
	Route() (method, path string)
	Middleware() []echo.MiddlewareFunc
	Serve(c echo.Context) error
}

func ProvideAs(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Router)),
		fx.ResultTags(`group:"global:http:router"`),
	)
}
