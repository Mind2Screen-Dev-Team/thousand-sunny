package http_middleware_global

import (
	"context"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/rs/xid"

	"github.com/labstack/echo/v4"
)

func ProvideTraceID() TraceID {
	return TraceID{}
}

type TraceID struct{}

func (TraceID) Name() string {
	return "req.trace.id"
}

func (TraceID) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			r   = c.Request()
			ctx = r.Context()
			id  = xid.New()
		)

		c.Set(xlog.XLOG_TRACE_ID_KEY, id)
		c.SetRequest(
			r.WithContext(context.WithValue(ctx, xlog.XLOG_REQ_TRACE_ID_CTX_KEY, id)),
		)

		return next(c)
	}
}
