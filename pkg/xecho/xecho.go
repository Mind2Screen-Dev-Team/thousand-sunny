package xecho

import (
	"github.com/ggicci/httpin"
	"github.com/labstack/echo/v4"
)

type (
	HttpInBinder struct{}
)

func (HttpInBinder) Bind(input any, c echo.Context) error {
	return httpin.DecodeTo(c.Request(), input)
}
