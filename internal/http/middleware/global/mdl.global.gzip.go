package http_middleware_global

import (
	"compress/gzip"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ProvideGZip() GZip {
	return GZip{}
}

type GZip struct{}

func (GZip) Name() string {
	return "gzip"
}

func (GZip) Order() int {
	return 5
}

func (GZip) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	cfg := middleware.GzipConfig{
		Skipper:   middleware.DefaultSkipper,
		Level:     gzip.DefaultCompression,
		MinLength: 0,
	}
	return middleware.GzipWithConfig(cfg)(next)
}
