package handler_user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xresp"

	http_middleware_private "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware/private"
	service_api "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service/api"
	service_attr "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service/attr"
)

type UserCreateParamFx struct {
	fx.In

	UserService service_api.UserServiceAPI
	AuthJWT     *http_middleware_private.AuthJWT
}

type UserCreateHandlerFx struct {
	UserCreateParamFx
}

func NewUserCreateHandlerFx(p UserCreateParamFx) UserCreateHandlerFx {
	return UserCreateHandlerFx{p}
}

func (UserCreateHandlerFx) Route() (method string, path string) {
	return http.MethodPost, "/api/v1/user"
}

func (u UserCreateHandlerFx) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		u.AuthJWT.Serve,
	}
}

func (u UserCreateHandlerFx) Serve(c echo.Context) error {
	var (
		r   = c.Request()
		ctx = r.Context()
		req = UserCreateDTO{}
		res = xresp.NewRestResponse[any, any](c)
	)

	if err := c.Bind(&req); err != nil {
		return res.
			StatusCode(http.StatusBadRequest).
			Code(http.StatusBadRequest).
			Msg("failed read request body").
			JSON()
	}

	var (
		body = service_attr.UserCreate{
			Name: req.Payload.Name,
			Age:  req.Payload.Age,
		}
	)

	if err := u.UserService.Create(ctx, body); err != nil {
		return res.
			StatusCode(http.StatusBadRequest).
			Code(http.StatusBadRequest).
			Msg("failed save user").
			JSON()
	}

	return res.
		StatusCode(http.StatusOK).
		Code(http.StatusOK).
		Msg("success create user").
		JSON()
}
