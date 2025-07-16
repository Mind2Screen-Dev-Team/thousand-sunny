package xfiber

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SkipPath(c *fiber.Ctx, ps ...string) (next func() error, skip bool) {
	p := c.Path()
	for _, v := range ps {
		if strings.HasPrefix(p, v) {
			return c.Next, true
		}
	}

	return nil, false
}
