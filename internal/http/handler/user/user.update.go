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
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service/user/attr"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhuma"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xresp"
)

var ExampleUserUpdateHandlerModuleFx = fx.Options(
	fx.Provide(xhuma.AnnotateHandlerAs(NewExampleUserUpdateHandlerFx)),
)

type ExampleUserUpdateHandlerParamFx struct {
	fx.In

	AuthJWT   *private.AuthJWT
	ExUserSvc api.ExampleUserServiceAPI
	LogDebug  *xlog.DebugLogger
}

type ExampleUserUpdateHandlerFx struct {
	p      ExampleUserUpdateHandlerParamFx
	logger xlog.Logger
}

func NewExampleUserUpdateHandlerFx(p ExampleUserUpdateHandlerParamFx) ExampleUserUpdateHandlerFx {
	return ExampleUserUpdateHandlerFx{p: p, logger: xlog.NewLogger(p.LogDebug.Logger)}
}

func (h ExampleUserUpdateHandlerFx) Register(api huma.API) {
	huma.Register(api, h.Operation(), h.Serve)
}

func (h ExampleUserUpdateHandlerFx) Operation() huma.Operation {
	return huma.Operation{
		OperationID:   "api-update-user",
		Path:          "/api/v1/user/{id}",
		Method:        http.MethodPut,
		Summary:       "Update User",
		Description:   "Updates an existing user's information based on the provided data. Returns the updated user's data or an error if the user is not found or the request is invalid.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Users"},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Description: "Successful response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Ref: "schemas/ExampleUserUpdateResponseBody",
						},
						Example: ExampleUserUpdateResponseBody{
							Code: http.StatusOK,
							Msg:  "ok",
							Data: &ExampleUserUpdateResponseData{
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
			strconv.Itoa(http.StatusUnprocessableEntity): {
				Description: "Validation failed response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Ref: "schemas/GeneralResponseError",
						},
						Example: xresp.GeneralResponseError{
							Code: http.StatusUnprocessableEntity,
							Msg:  http.StatusText(http.StatusUnprocessableEntity),
							Data: nil,
							Err: &xresp.ErrorModel{
								Detail: "validation failed",
								Errors: []*huma.ErrorDetail{
									{
										Message:  "expected number >= 18",
										Location: "body.age",
										Value:    11,
									},
									{
										Message:  "expected length >= 5",
										Location: "body.name",
										Value:    "John",
									},
								},
							},
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
						Example: ExampleUserUpdateResponseBody{
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

func (h ExampleUserUpdateHandlerFx) Serve(ctx context.Context, in *ExampleUserUpdateRequestInput) (out *ExampleUserUpdateResponseOutput, err error) {
	d, err := h.p.ExUserSvc.Update(ctx, attr.ExampleUser{
		ID:   in.ID,
		Name: in.Body.Name,
		Age:  in.Body.Age,
	})
	if err != nil {
		h.logger.Error(ctx, "failed to update user", "input", in, "err", fmt.Sprintf("%+v", err))
		return nil, huma.Error500InternalServerError("failed to update user", err)
	}

	var (
		body = ExampleUserUpdateResponseBody{
			Code: http.StatusOK,
			Msg:  "ok",
			Data: &ExampleUserUpdateResponseData{
				ID:        d.ID,
				Name:      d.Name,
				Age:       d.Age,
				CreatedAt: d.UpdatedAt,
				UpdatedAt: d.UpdatedAt,
			},
		}

		resp = ExampleUserUpdateResponseOutput{
			Status: http.StatusOK,
			Body:   body,
		}
	)

	return &resp, nil
}
