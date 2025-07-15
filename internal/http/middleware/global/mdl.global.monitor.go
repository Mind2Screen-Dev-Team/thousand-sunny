package global

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func ProvideMonitor() Monitor {
	return Monitor{}
}

type Monitor struct{}

func (Monitor) Name() string {
	return "monitor"
}

func (Monitor) App(app *fiber.App) {
	app.Get("/monitor", monitor.New(monitor.ConfigDefault))
}

func (Monitor) Serve(c *fiber.Ctx) error {
	return c.Next()
}
