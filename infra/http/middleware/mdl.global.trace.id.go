package middleware

import (
	"context"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/constant"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xfiber"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xtracer"
	"github.com/gofiber/fiber/v2"

	"github.com/rs/xid"

	"go.opentelemetry.io/otel/trace"
)

func ProvideTraceID(tracer trace.Tracer) TraceID {
	return TraceID{tracer}
}

type TraceID struct {
	tracer trace.Tracer
}

func (TraceID) App(app *fiber.App) {}

func (TraceID) Name() string {
	return "trace.id"
}

func (t TraceID) Serve(c *fiber.Ctx) error {
	if next, ok := xfiber.SkipPath(c, constant.FiberSkipablePathFromMiddleware[:]...); ok {
		return next()
	}

	var (
		ctx  = c.UserContext()
		id   = xid.New()
		span trace.Span
	)

	ctx, span = xtracer.Start(t.tracer, ctx, "global trace id")
	defer span.End()

	c.SetUserContext(context.WithValue(ctx, xlog.XLOG_REQ_TRACE_ID_CTX_KEY, id.String()))

	return c.Next()
}
