package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/http/middleware/private"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service/user/api"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhuma"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
)

var ExampleUserReadAllHandlerModuleFx = fx.Options(
	fx.Provide(xhuma.AnnotateHandlerAs(NewExampleUserReadAllHandlerFx)),
)

type ExampleUserReadAllHandlerParamFx struct {
	fx.In

	AuthJWT   *private.AuthJWT
	ExUserSvc api.ExampleUserServiceAPI
	LogDebug  *xlog.DebugLogger
}

type ExampleUserReadAllHandlerFx struct {
	p      ExampleUserReadAllHandlerParamFx
	logger xlog.Logger
}

func NewExampleUserReadAllHandlerFx(p ExampleUserReadAllHandlerParamFx) ExampleUserReadAllHandlerFx {
	return ExampleUserReadAllHandlerFx{p: p, logger: xlog.NewLogger(p.LogDebug.Logger)}
}

func (h ExampleUserReadAllHandlerFx) Register(api huma.API) {
	huma.Register(api, h.Operation(), h.Serve)
}

func (h ExampleUserReadAllHandlerFx) Operation() huma.Operation {
	return huma.Operation{
		OperationID:   "api-read-all-user",
		Path:          "/api/v1/users",
		Method:        http.MethodGet,
		Summary:       "Retrieves All Users",
		Description:   "Retrieves a list of all users.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Users"},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Description: "Successful response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Ref: "schemas/ExampleUserReadAllResponseBody",
						},
						Example: ExampleUserReadAllResponseBody{
							Code: http.StatusOK,
							Msg:  "ok",
							Data: []ExampleUserReadAllResponseData{
								{
									ID:        uuid.Must(uuid.NewV7()),
									Name:      "Johnny",
									Age:       18,
									CreatedAt: time.Now(),
									UpdatedAt: time.Now().Add(5 * time.Hour),
								},
								{
									ID:        uuid.Must(uuid.NewV7()),
									Name:      "Rebecca",
									Age:       19,
									CreatedAt: time.Now(),
									UpdatedAt: time.Now().Add(5 * time.Hour),
								},
							},
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
						Example: ExampleUserReadAllResponseBody{
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

func (h ExampleUserReadAllHandlerFx) Serve(ctx context.Context, in *ExampleUserReadAllRequestInput) (out *ExampleUserReadAllResponseOutput, err error) {
	d, err := h.p.ExUserSvc.ReadAll(ctx, 100, 0)
	if err != nil {
		h.logger.Error(ctx, "failed to read all user", "input", in, "err", fmt.Sprintf("%+v", err))
	}

	dd := make([]ExampleUserReadAllResponseData, len(d))
	for i, v := range d {
		dd[i] = ExampleUserReadAllResponseData{
			ID:        v.ID,
			Name:      v.Name,
			Age:       v.Age,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
	}

	var (
		body = ExampleUserReadAllResponseBody{
			Code: http.StatusOK,
			Msg:  "ok",
			Data: dd,
		}

		resp = ExampleUserReadAllResponseOutput{
			Status: http.StatusOK,
			Body:   body,
		}
	)

	return &resp, nil
}
