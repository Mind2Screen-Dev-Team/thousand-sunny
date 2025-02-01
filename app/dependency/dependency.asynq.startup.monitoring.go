package dependency

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hibiken/asynqmon"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

type InvokeAsynqMonitoringServerParam struct {
	fx.In

	Echo *echo.Echo
	Opt  *asynqmon.Options
}

func InvokeAsynqMonitoringServer(p InvokeAsynqMonitoringServerParam) error {
	if p.Echo == nil {
		return errors.New("field 'Echo' with type '*echo.Echo' is not provided")
	}

	if p.Opt == nil {
		return errors.New("field 'Opt' with type '*asynqmon.Options' is not provided")
	}

	var (
		_asynqmon = asynqmon.New(*p.Opt)
		_rootpath = fmt.Sprintf("%s/*", p.Opt.RootPath)
	)

	p.Echo.Use(middleware.CORS())
	p.Echo.Use(middleware.Recover())
	p.Echo.Any(_rootpath, echo.WrapHandler(_asynqmon))

	p.Echo.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, ".")
	})

	p.Echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, http.StatusText(http.StatusOK))
	})
	return nil
}
