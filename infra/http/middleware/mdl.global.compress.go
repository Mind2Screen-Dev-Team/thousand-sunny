package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

func ProvideCompress() Compress {
	return Compress{}
}

type Compress struct{}

func (Compress) Name() string {
	return "compress"
}

func (Compress) App(app *fiber.App) {}

func (Compress) Serve(c *fiber.Ctx) error {
	var (
		cfg = compress.Config{
			Level: compress.LevelBestSpeed,
		}
	)

	return compress.New(cfg)(c)
}
