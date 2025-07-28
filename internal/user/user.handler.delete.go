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

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
)

type ExampleUserDeleteHandlerParamFx struct {
	fx.In

	ExUserSvc ExampleUserServiceAPI
	LogDebug  *xlog.DebugLogger
}

type ExampleUserDeleteHandlerFx struct {
	p      ExampleUserDeleteHandlerParamFx
	logger xlog.Logger
}

func NewDeleteHandlerFx(p ExampleUserDeleteHandlerParamFx) ExampleUserDeleteHandlerFx {
	return ExampleUserDeleteHandlerFx{p: p, logger: xlog.NewLogger(p.LogDebug.Logger)}
}

func (h ExampleUserDeleteHandlerFx) Register(api huma.API) {
	huma.Register(api, h.Operation(), h.Serve)
}

func (h ExampleUserDeleteHandlerFx) Operation() huma.Operation {
	return huma.Operation{
		OperationID:   "api-delete-user",
		Path:          "/api/v1/user/{id}",
		Method:        http.MethodDelete,
		Summary:       "Delete User",
		Description:   "Deletes a specific user identified by their unique ID. Returns a success status if the deletion is successful, or an error if the user does not exist.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Users"},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Description: "Successful response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Ref: "schemas/ExampleUserDeleteResponseBody",
						},
						Example: ExampleUserDeleteResponseBody{
							Code: http.StatusOK,
							Msg:  "ok",
							Data: &ExampleUserDeleteResponseData{
								ID:        uuid.Must(uuid.NewV7()),
								Name:      "Johnny",
								Age:       18,
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
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
						Example: ExampleUserDeleteResponseBody{
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
		
func (h ExampleUserDeleteHandlerFx) Serve(ctx context.Context, in *ExampleUserDeleteRequestInput) (out *ExampleUserDeleteResponseOutput, err error) {
	d, err := h.p.ExUserSvc.Delete(ctx, in.ID.String())
	if err != nil {
		h.logger.Error(ctx, "failed to delete user", "input", in, "err", fmt.Sprintf("%+v", err))
		return nil, huma.Error500InternalServerError("failed to delete user", err)
	}

	var (
		body = ExampleUserDeleteResponseBody{
			Code: http.StatusOK,
			Msg:  "ok",
			Data: &ExampleUserDeleteResponseData{
				ID:        d.ID,
				Name:      d.Name,
				Age:       d.Age,
				CreatedAt: d.CreatedAt,
				UpdatedAt: d.UpdatedAt,
			},
		}

		resp = ExampleUserDeleteResponseOutput{
			Status: http.StatusOK,
			Body:   body,
		}
	)

	return &resp, nil
}
