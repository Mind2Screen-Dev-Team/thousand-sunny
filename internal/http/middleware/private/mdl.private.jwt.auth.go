package http_middleware_private

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type (
	AuthJWTParams struct {
		fx.In

		Cfg   config.Cfg
		Debug *xlog.DebugLogger
	}

	AuthJWT struct {
		cfg   config.Cfg
		debug xlog.Logger
	}
)

func NewAuthJWT(p AuthJWTParams) (*AuthJWT, error) {
	if p.Debug == nil {
		return nil, errors.New("field 'Debug' with type '*xlog.DebugLogger' is not provided")
	}

	return &AuthJWT{cfg: p.Cfg, debug: xlog.NewLogger(p.Debug.Logger)}, nil
}

func (a AuthJWT) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			req  = c.Request()
			ctx  = req.Context()
			auth = req.Header.Get("Authorization")
			msg  = http.StatusText(http.StatusUnauthorized)
		)

		if !strings.HasPrefix(auth, "Bearer ") {
			return c.String(http.StatusUnauthorized, msg)
		}

		if token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer ")); token != "abc" {
			return c.String(http.StatusUnauthorized, msg)
		}

		a.debug.Info(ctx, "auth is success")

		return next(c)
	}
}
