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

var ExampleUserReadHandlerModuleFx = fx.Options(
	fx.Provide(xhuma.AnnotateHandlerAs(NewExampleUserReadHandlerFx)),
)

type ExampleUserReadHandlerParamFx struct {
	fx.In

	AuthJWT   *private.AuthJWT
	ExUserSvc api.ExampleUserServiceAPI
	LogDebug  *xlog.DebugLogger
}

type ExampleUserReadHandlerFx struct {
	p      ExampleUserReadHandlerParamFx
	logger xlog.Logger
}

func NewExampleUserReadHandlerFx(p ExampleUserReadHandlerParamFx) ExampleUserReadHandlerFx {
	return ExampleUserReadHandlerFx{p: p, logger: xlog.NewLogger(p.LogDebug.Logger)}
}

func (h ExampleUserReadHandlerFx) Register(api huma.API) {
	huma.Register(api, h.Operation(), h.Serve)
}

func (h ExampleUserReadHandlerFx) Operation() huma.Operation {
	return huma.Operation{
		OperationID:   "api-read-user",
		Path:          "/api/v1/user/{id}",
		Method:        http.MethodGet,
		Summary:       "Retrieves User Details",
		Description:   "Retrieves detailed information about a specific user identified by their unique ID. Returns an error if the user does not exist.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Users"},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Description: "Successful response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Ref: "schemas/ExampleUserReadResponseBody",
						},
						Example: ExampleUserReadResponseBody{
							Code: http.StatusOK,
							Msg:  "ok",
							Data: &ExampleUserReadResponseData{
								ID:        uuid.Must(uuid.NewV7()),
								Name:      "Johnny",
								Age:       18,
								CreatedAt: time.Now(),
								UpdatedAt: time.Now().Add(5 * time.Hour),
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
						Example: ExampleUserReadResponseBody{
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

func (h ExampleUserReadHandlerFx) Serve(ctx context.Context, in *ExampleUserReadRequestInput) (out *ExampleUserReadResponseOutput, err error) {
	d, err := h.p.ExUserSvc.Read(ctx, in.ID.String())
	if err != nil {
		h.logger.Error(ctx, "failed to read detail user", "input", in, "err", fmt.Sprintf("%+v", err))
		return nil, huma.Error500InternalServerError("failed to read detail user", err)
	}

	var (
		body = ExampleUserReadResponseBody{
			Code: http.StatusOK,
			Msg:  "ok",
			Data: &ExampleUserReadResponseData{
				ID:        d.ID,
				Name:      d.Name,
				Age:       d.Age,
				CreatedAt: d.CreatedAt,
				UpdatedAt: d.UpdatedAt,
			},
		}

		resp = ExampleUserReadResponseOutput{
			Status: http.StatusOK,
			Body:   body,
		}
	)

	return &resp, nil
}
