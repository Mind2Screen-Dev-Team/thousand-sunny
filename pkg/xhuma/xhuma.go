package xhuma

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xresp"
)

// Repalce Existing Huma New Error
func init() {
	huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
		details := make([]*huma.ErrorDetail, len(errs))
		for i := range errs {
			if converted, ok := errs[i].(huma.ErrorDetailer); ok {
				details[i] = converted.ErrorDetail()
			} else {
				if errs[i] == nil {
					continue
				}
				details[i] = &huma.ErrorDetail{Message: errs[i].Error()}
			}
		}
		return &xresp.GeneralResponseError{
			Code: status,
			Msg:  http.StatusText(status),
			Err: &xresp.ErrorModel{
				Detail: msg,
				Errors: details,
			},
		}
	}

	huma.NewErrorWithContext = func(_ huma.Context, status int, msg string, errs ...error) huma.StatusError {
		return huma.NewError(status, msg, errs...)
	}
}

type HandlerRegister interface {
	Register(api huma.API)
}

func AnnotateHandlerAs(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(HandlerRegister)),
		fx.ResultTags(`group:"global:http:handler"`),
	)
}

type GlobalMiddleware interface {
	Name() string
	App(app *fiber.App)
	Serve(c *fiber.Ctx) error
}

func AnnotateGlobalMiddlewareAs(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(GlobalMiddleware)),
		fx.ResultTags(`group:"global:http:middleware"`),
	)
}
