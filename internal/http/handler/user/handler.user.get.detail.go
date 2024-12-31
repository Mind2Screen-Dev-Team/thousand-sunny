package handler_user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xresp"

	http_middleware_private "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware/private"
	service_api "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service/api"
)

type UserGetDetailParamFx struct {
	fx.In

	UserService service_api.UserServiceAPI
	AuthJWT     *http_middleware_private.AuthJWT
}

type UserGetDetailHandlerFx struct {
	UserGetDetailParamFx
}

func NewUserGetDetailHandlerFx(p UserGetDetailParamFx) UserGetDetailHandlerFx {
	return UserGetDetailHandlerFx{p}
}

func (UserGetDetailHandlerFx) Route() (method string, path string) {
	return http.MethodGet, "/api/v1/user"
}

func (u UserGetDetailHandlerFx) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		u.AuthJWT.Serve,
	}
}

func (UserGetDetailHandlerFx) Serve(c echo.Context) error {
	var (
		res = xresp.NewRestResponse[any, any](c)
	)

	return res.
		StatusCode(http.StatusOK).
		Code(http.StatusOK).
		Msg("success get user").
		Data(map[string]any{"status": true}).
		JSON()
}
