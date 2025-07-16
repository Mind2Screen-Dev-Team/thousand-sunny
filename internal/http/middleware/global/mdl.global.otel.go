package global

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/constant"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xfiber"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
)

func ProvideOtel() Otel {
	return Otel{}
}

type Otel struct{}

func (Otel) Name() string {
	return "otel.http"
}

func (Otel) App(app *fiber.App) {}

func (m Otel) Serve(c *fiber.Ctx) error {
	if next, ok := xfiber.SkipPath(c, constant.FiberSkipablePathFromMiddleware[:]...); ok {
		return next()
	}

	return otelfiber.Middleware(
		otelfiber.WithCollectClientIP(true),
	)(c)
}
