package global

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
)

func ProvideFavicon() Favicon {
	return Favicon{}
}

type Favicon struct{}

func (Favicon) Name() string {
	return "monitor"
}

func (Favicon) App(app *fiber.App) {}

func (Favicon) Serve(c *fiber.Ctx) error {
	return favicon.New(favicon.Config{
		File:         "storage/assets/favicon.ico",
		URL:          "/favicon.ico",
		CacheControl: "public, max-age=31536000",
	})(c)
}
