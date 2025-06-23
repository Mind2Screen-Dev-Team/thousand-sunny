package xecho

import (
	"context"
	"encoding/json"

	"github.com/ggicci/httpin"
	"github.com/labstack/echo/v4"
)

type (
	CtxTypeEnv   struct{}
	CtxUser      struct{}
	HttpInBinder struct{}
)

func (HttpInBinder) Bind(input any, c echo.Context) error {
	return httpin.DecodeTo(c.Request(), input)
}

func SetTypeEnv(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, CtxTypeEnv{}, value)
}

func GetTypeEnv(ctx context.Context) string {
	value, ok := ctx.Value(CtxTypeEnv{}).(string)
	if !ok {
		value = "development"
	}
	return value
}

func SetUser(ctx context.Context, value any) context.Context {
	data, err := json.Marshal(value)
	if err != nil {
		data = []byte("")
	}
	return context.WithValue(ctx, CtxUser{}, string(data))
}

func GetUser(ctx context.Context) string {
	value, ok := ctx.Value(CtxUser{}).(string)
	if !ok {
		value = ""
	}
	return value
}

func GetParsedUser(ctx context.Context, p any) error {
	value, ok := ctx.Value(CtxUser{}).(string)
	if !ok {
		value = ""
	}
	return json.Unmarshal([]byte(value), p)
}
