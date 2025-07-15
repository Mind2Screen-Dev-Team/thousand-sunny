package global

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ProvideCORS() CORS {
	return CORS{}
}

type CORS struct{}

func (CORS) Name() string {
	return "cors"
}

func (CORS) App(app *fiber.App) {}

func (CORS) Serve(c *fiber.Ctx) error {
	return cors.New()(c)
}
