package http_middleware

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ProvideRequestID() RequestID {
	return RequestID{}
}

type RequestID struct{}

func (RequestID) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(xlog.XLOG_TRACE_ID_KEY, uuid.Must(uuid.NewV7()))
		return next(c)
	}
}
