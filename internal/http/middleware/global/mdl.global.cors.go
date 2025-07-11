package http_middleware_global

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ProvideCORS() CORS {
	return CORS{}
}

type CORS struct{}

func (CORS) Name() string {
	return "cors"
}

func (CORS) Order() int {
	return 1
}

func (CORS) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.CORS()(next)
}
