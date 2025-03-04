package xlog

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

type CtxKey string

func (s CtxKey) String() string {
	return string(s)
}

const (
	XLOG_REQ_TRACE_ID_CTX_KEY CtxKey = "XLOG_REQ_TRACE_ID_CTX_KEY"
)

const (
	XLOG_TRACE_ID_KEY string = "XLOG_TRACE_ID_KEY"
	XLOG_KEY          string = "XLOG_KEY"
)

var (
	nopZeroLogger = zerolog.Nop()
)

func FromEcho(c echo.Context) Logger {
	v, ok := c.Get(XLOG_KEY).(zerolog.Logger)
	if !ok {
		v = nopZeroLogger
	}

	return NewLogger(
		v.With().Any("req_trace_id", c.Get(XLOG_TRACE_ID_KEY)).Logger(),
	)
}

type Logger interface {
	/*
		Fields is a helper function to use a map or slice to set fields using type assertion.
		[]any must alternate string keys and arbitrary values, and extraneous ones are ignored. i.e:

		With Request HTTP Context:

			var (
				first  = "first value"
				second = "second value"
			)
			<Imported_Log>.Trace(ctx, "hello", "first", first, "second", second)
			<Imported_Log>.Trace(ctx, xlog.Msgf("hello %s", "world!"), "first", first, "second", second)
			<Imported_Log>.Trace(ctx, "oh snap! got error", "error", err)
	*/
	Trace(ctx context.Context, msg string, fields ...any)

	/*
		Fields is a helper function to use a map or slice to set fields using type assertion.
		[]any must alternate string keys and arbitrary values, and extraneous ones are ignored. i.e:

		With Request HTTP Context:

			var (
				first  = "first value"
				second = "second value"
			)
			<Imported_Log>.Debug(ctx, "hello", "first", first, "second", second)
			<Imported_Log>.Debug(ctx, xlog.Msgf("hello %s", "world!"), "first", first, "second", second)
			<Imported_Log>.Debug(ctx, "oh snap! got error", "error", err)
	*/
	Debug(ctx context.Context, msg string, fields ...any)

	/*
		Fields is a helper function to use a map or slice to set fields using type assertion.
		[]any must alternate string keys and arbitrary values, and extraneous ones are ignored. i.e:

		With Request HTTP Context:

			var (
				first  = "first value"
				second = "second value"
			)
			<Imported_Log>.Info(ctx, "hello", "first", first, "second", second)
			<Imported_Log>.Info(ctx, xlog.Msgf("hello %s", "world!"), "first", first, "second", second)
			<Imported_Log>.Info(ctx, "oh snap! got error", "error", err)
	*/
	Info(ctx context.Context, msg string, fields ...any)

	/*
		Fields is a helper function to use a map or slice to set fields using type assertion.
		[]any must alternate string keys and arbitrary values, and extraneous ones are ignored. i.e:

		With Request HTTP Context:

			var (
				first  = "first value"
				second = "second value"
			)
			<Imported_Log>.Warn(ctx, "hello", "first", first, "second", second)
			<Imported_Log>.Warn(ctx, xlog.Msgf("hello %s", "world!"), "first", first, "second", second)
			<Imported_Log>.Warn(ctx, "oh snap! got error", "error", err)
	*/
	Warn(ctx context.Context, msg string, fields ...any)

	/*
		Fields is a helper function to use a map or slice to set fields using type assertion.
		[]any must alternate string keys and arbitrary values, and extraneous ones are ignored. i.e:

		With Request HTTP Context:

			var (
				first  = "first value"
				second = "second value"
			)
			<Imported_Log>.Error(ctx, "hello", "first", first, "second", second)
			<Imported_Log>.Error(ctx, xlog.Msgf("hello %s", "world!"), "first", first, "second", second)
			<Imported_Log>.Error(ctx, "oh snap! got error", "error", err)
	*/
	Error(ctx context.Context, msg string, fields ...any)

	/*
		Fields is a helper function to use a map or slice to set fields using type assertion.
		[]any must alternate string keys and arbitrary values, and extraneous ones are ignored. i.e:

		With Request HTTP Context:

			var (
				first  = "first value"
				second = "second value"
			)
			<Imported_Log>.Fatal(ctx, "hello", "first", first, "second", second)
			<Imported_Log>.Fatal(ctx, xlog.Msgf("hello %s", "world!"), "first", first, "second", second)
			<Imported_Log>.Fatal(ctx, "oh snap! got error", "error", err)
	*/
	Fatal(ctx context.Context, msg string, fields ...any)

	/*
		Fields is a helper function to use a map or slice to set fields using type assertion.
		[]any must alternate string keys and arbitrary values, and extraneous ones are ignored. i.e:

		With Request HTTP Context:

			var (
				first  = "first value"
				second = "second value"
			)
			<Imported_Log>.Panic(ctx, "hello", "first", first, "second", second)
			<Imported_Log>.Panic(ctx, xlog.Msgf("hello %s", "world!"), "first", first, "second", second)
			<Imported_Log>.Panic(ctx, "oh snap! got error", "error", err)
	*/
	Panic(ctx context.Context, msg string, fields ...any)
}

func Msgf(ctx context.Context, msg string, args ...any) string {
	return fmt.Sprintf(msg, args...)
}

type ZeroLogger struct {
	log zerolog.Logger
}

func NewLogger(log zerolog.Logger) Logger {
	return &ZeroLogger{log}
}

// attachFields directly to the zerolog.Event object without creating a new map
func (zl *ZeroLogger) attachFields(ctx context.Context, e *zerolog.Event, fields []any) *zerolog.Event {
	span := trace.SpanFromContext(ctx).SpanContext()
	if span.IsValid() {
		fields = append(fields,
			"otel_span_id", span.SpanID().String(),
			"otel_trace_id", span.TraceID().String(),
		)
	}

	if v, ok := ctx.Value(XLOG_REQ_TRACE_ID_CTX_KEY).(xid.ID); ok {
		fields = append(fields, "req_trace_id", v)
	}

	for i := 0; i < len(fields); i += 2 {
		if isHasKV := i+1 < len(fields); !isHasKV {
			continue
		}

		key, ok := fields[i].(string)
		if !ok {
			continue
		}

		e = AnyFieldToZeroLogEvent(e, key, fields[i+1])
	}

	return e
}

func (zl *ZeroLogger) Trace(ctx context.Context, msg string, fields ...any) {
	e := zl.log.Trace()
	e = e.Ctx(ctx)
	e = zl.attachFields(ctx, e, fields)
	e.Msg(msg)
}

func (zl *ZeroLogger) Debug(ctx context.Context, msg string, fields ...any) {
	e := zl.log.Debug()
	e = e.Ctx(ctx)
	e = zl.attachFields(ctx, e, fields)
	e.Msg(msg)
}

func (zl *ZeroLogger) Info(ctx context.Context, msg string, fields ...any) {
	e := zl.log.Info()
	e = e.Ctx(ctx)
	e = zl.attachFields(ctx, e, fields)
	e.Msg(msg)
}

func (zl *ZeroLogger) Warn(ctx context.Context, msg string, fields ...any) {
	e := zl.log.Warn()
	e = e.Ctx(ctx)
	e = zl.attachFields(ctx, e, fields)
	e.Msg(msg)
}

func (zl *ZeroLogger) Error(ctx context.Context, msg string, fields ...any) {
	e := zl.log.Error()
	e = e.Ctx(ctx)
	e = zl.attachFields(ctx, e, fields)
	e.Msg(msg)
}

func (zl *ZeroLogger) Fatal(ctx context.Context, msg string, fields ...any) {
	e := zl.log.Fatal()
	e = e.Ctx(ctx)
	e = zl.attachFields(ctx, e, fields)
	e.Msg(msg)
}

func (zl *ZeroLogger) Panic(ctx context.Context, msg string, fields ...any) {
	e := zl.log.Panic()
	e = e.Ctx(ctx)
	e = zl.attachFields(ctx, e, fields)
	e.Msg(msg)
}
