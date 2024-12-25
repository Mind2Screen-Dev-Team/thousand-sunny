package http_router

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Middleware interface {
	Serve(next echo.HandlerFunc) echo.HandlerFunc
}

func As(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Middleware)),
		fx.ResultTags(`group:"routers"`),
	)
}
