package handler_health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct{}

func NewHealthHandler() HealthHandler {
	return HealthHandler{}
}

func (HealthHandler) Route() (method, path string) {
	return http.MethodGet, "/api/v1/health"
}

func (u HealthHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (HealthHandler) Serve(c echo.Context) error {
	return c.String(http.StatusOK, ".")
}
