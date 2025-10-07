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

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhuma"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xresp"
)

type ExampleUserCreateHandlerParamFx struct {
	fx.In

	ExUserSvc ExampleUserServiceAPI
	LogDebug  *xlog.DebugLogger
}

type ExampleUserCreateHandlerFx struct {
	p      ExampleUserCreateHandlerParamFx
	logger xlog.Logger
}

type ExampleUserCreateHandlerFxOut struct {
	fx.Out

	Handler xhuma.HandlerRegister `group:"global:http:handler"`
}

func NewCreateHandlerFx(p ExampleUserCreateHandlerParamFx) ExampleUserCreateHandlerFxOut {
	return ExampleUserCreateHandlerFxOut{
		Handler: &ExampleUserCreateHandlerFx{p: p, logger: xlog.NewLogger(p.LogDebug.Logger)},
	}
}

func (h ExampleUserCreateHandlerFx) Register(api huma.API) {
	huma.Register(api, h.Operation(), h.Serve)
}

func (h ExampleUserCreateHandlerFx) Operation() huma.Operation {
	return huma.Operation{
		OperationID:   "api-create-user",
		Path:          "/api/v1/user",
		Method:        http.MethodPost,
		Summary:       "Create New User",
		Description:   "Creates a new user with the provided information and returns the created user's data or an error.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Users"},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Description: "Successful response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Ref: "schemas/ExampleUserCreateResponseBody",
						},
						Example: ExampleUserCreateResponseBody{
							Code: http.StatusOK,
							Msg:  "ok",
							Data: &ExampleUserCreateResponseData{
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
						Example: ExampleUserCreateResponseBody{
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

func (h ExampleUserCreateHandlerFx) Serve(ctx context.Context, in *ExampleUserCreateRequestInput) (out *ExampleUserCreateResponseOutput, err error) {
	d, err := h.p.ExUserSvc.Create(ctx, ExampleUser{
		Name: in.Body.Name,
		Age:  in.Body.Age,
	})
	if err != nil {
		h.logger.Error(ctx, "failed to create new user", "input", in, "err", fmt.Sprintf("%+v", err))
		return nil, huma.Error500InternalServerError("failed to create new user", err)
	}

	var (
		body = ExampleUserCreateResponseBody{
			Code: http.StatusOK,
			Msg:  "ok",
			Data: &ExampleUserCreateResponseData{
				ID:        d.ID,
				Name:      d.Name,
				Age:       d.Age,
				CreatedAt: d.CreatedAt,
				UpdatedAt: d.UpdatedAt,
			},
		}

		resp = ExampleUserCreateResponseOutput{
			Status: http.StatusOK,
			Body:   body,
		}
	)

	return &resp, nil
}
