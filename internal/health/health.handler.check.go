package health

import (
	"context"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/infra/http/middleware"
)

type HealthHandlerParamFx struct {
	fx.In

	PrivateAuthJWT *middleware.PrivateAuthJWT
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
