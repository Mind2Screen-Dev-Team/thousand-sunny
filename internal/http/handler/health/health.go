package health

import (
	"context"
	"net/http"
	"strconv"

	http_middleware_private "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware/private"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhuma"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/xid"
	"go.uber.org/fx"
)

var HealthHandlerModuleFx = fx.Options(
	fx.Provide(xhuma.AnnotateHandlerAs(NewHealthHandlerFx)),
)

type HealthHandlerParamFx struct {
	fx.In

	AuthJWT *http_middleware_private.AuthJWT
}

type HealthHandlerFx struct {
	p HealthHandlerParamFx
}

func NewHealthHandlerFx(p HealthHandlerParamFx) HealthHandlerFx {
	return HealthHandlerFx{p}
}

func (h HealthHandlerFx) Register(api huma.API) {
	huma.Register(api, h.Operation(), h.Serve)
}

func (h HealthHandlerFx) Operation() huma.Operation {
	return huma.Operation{
		OperationID:   "api-check-health",
		Path:          "/api/v1/health",
		Method:        http.MethodGet,
		Summary:       "Service Health",
		Description:   "Returns status ok if the service is healthy",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Liveness"},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Description: "Successful response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Ref: "schemas/HealthResponseBody",
						},
						Example: HealthResponseBody{
							Code:    http.StatusOK,
							Msg:     "ok",
							Data:    nil,
							Err:     nil,
							TraceID: xid.New().String(),
						},
					},
				},
			},
			strconv.Itoa(http.StatusInternalServerError): {
				Description: "Failed response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Ref: "schemas/GeneralResponseError",
						},
						Example: HealthResponseBody{
							Code:    http.StatusInternalServerError,
							Msg:     http.StatusText(http.StatusInternalServerError),
							Data:    nil,
							Err:     nil,
							TraceID: xid.New().String(),
						},
					},
				},
			},
		},
	}
}

func (h HealthHandlerFx) Serve(ctx context.Context, in *HealthRequestInput) (out *HealthResponseOutput, err error) {
	var (
		body = HealthResponseBody{
			Code: http.StatusOK,
			Msg:  "ok",
		}
		resp = HealthResponseOutput{
			Status: http.StatusOK,
			Body:   body,
		}
	)
	return &resp, nil
}
