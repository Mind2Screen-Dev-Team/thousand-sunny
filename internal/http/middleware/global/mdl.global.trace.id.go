package http_middleware_global

import (
	"context"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xtracer"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"

	"go.opentelemetry.io/otel/trace"
)

func ProvideTraceID(tracer trace.Tracer) TraceID {
	return TraceID{tracer}
}

type TraceID struct {
	tracer trace.Tracer
}

func (TraceID) Name() string {
	return "trace.id"
}

func (TraceID) Order() int {
	return 2
}

func (t TraceID) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			r    = c.Request()
			ctx  = r.Context()
			id   = xid.New()
			span trace.Span
		)

		ctx, span = xtracer.Start(t.tracer, ctx, "set internal trace id")
		defer span.End()

		c.Set(xlog.XLOG_TRACE_ID_KEY, id)
		c.SetRequest(
			r.WithContext(context.WithValue(ctx, xlog.XLOG_REQ_TRACE_ID_CTX_KEY, id)),
		)

		if err := next(c); err != nil {
			span.RecordError(err)
			return err
		}

		return nil
	}
}
