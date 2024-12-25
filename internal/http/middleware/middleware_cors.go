package http_middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ProvideCORS() CORS {
	return CORS{}
}

type CORS struct{}

func (CORS) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.CORS()(next)
}
