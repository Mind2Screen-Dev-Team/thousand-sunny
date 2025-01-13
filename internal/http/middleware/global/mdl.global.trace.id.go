package http_middleware_global

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/rs/xid"

	"github.com/labstack/echo/v4"
)

func ProvideTraceID() TraceID {
	return TraceID{}
}

type TraceID struct{}

func (TraceID) Name() string {
	return "trace.id"
}

func (TraceID) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(xlog.XLOG_TRACE_ID_KEY, xid.New())
		return next(c)
	}
}
