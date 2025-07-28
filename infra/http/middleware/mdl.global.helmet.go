package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

func ProvideHelmet() Helmet {
	return Helmet{}
}

type Helmet struct{}

func (Helmet) Name() string {
	return "helmet"
}

func (Helmet) App(app *fiber.App) {}

func (Helmet) Serve(c *fiber.Ctx) error {
	return helmet.New()(c)
}
