package xlog

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

const (
	XLOG_TRACE_ID_KEY = "XLOG_TRACE_ID_KEY"
	XLOG_KEY          = "XLOG_KEY"
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
		v.With().Any("trace_id", c.Get(XLOG_TRACE_ID_KEY)).Logger(),
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
			xlog.FromEcho(c).Trace("hello", "first", first, "second", second)
			xlog.FromEcho(c).Trace(xlogger.Msgf("hello %s", "world!"), "first", first, "second", second)
			xlog.FromEcho(c).Trace("oh snap! got error", "error", err)
	*/
	Trace(msg string, fields ...any)

	/*
		Fields is a helper function to use a map or slice to set fields using type assertion.
		[]any must alternate string keys and arbitrary values, and extraneous ones are ignored. i.e:

		With Request HTTP Context:

			var (
				first  = "first value"
				second = "second value"
			)
			xlog.FromEcho(c).Debug("hello", "first", first, "second", second)
			xlog.FromEcho(c).Debug(xlogger.Msgf("hello %s", "world!"), "first", first, "second", second)
			xlog.FromEcho(c).Debug("oh snap! got error", "error", err)
	*/
	Debug(msg string, fields ...any)

	/*
		Fields is a helper function to use a map or slice to set fields using type assertion.
		[]any must alternate string keys and arbitrary values, and extraneous ones are ignored. i.e:

		With Request HTTP Context:

			var (
				first  = "first value"
				second = "second value"
			)
			xlog.FromEcho(c).Info("hello", "first", first, "second", second)
			xlog.FromEcho(c).Info(xlogger.Msgf("hello %s", "world!"), "first", first, "second", second)
			xlog.FromEcho(c).Info("oh snap! got error", "error", err)
	*/
	Info(msg string, fields ...any)

	/*
		Fields is a helper function to use a map or slice to set fields using type assertion.
		[]any must alternate string keys and arbitrary values, and extraneous ones are ignored. i.e:

		With Request HTTP Context:

			var (
				first  = "first value"
				second = "second value"
			)
			xlog.FromEcho(c).Warn("hello", "first", first, "second", second)
			xlog.FromEcho(c).Warn(xlogger.Msgf("hello %s", "world!"), "first", first, "second", second)
			xlog.FromEcho(c).Warn("oh snap! got error", "error", err)
	*/
	Warn(msg string, fields ...any)

	/*
		Fields is a helper function to use a map or slice to set fields using type assertion.
		[]any must alternate string keys and arbitrary values, and extraneous ones are ignored. i.e:

		With Request HTTP Context:

			var (
				first  = "first value"
				second = "second value"
			)
			xlog.FromEcho(c).Error("hello", "first", first, "second", second)
			xlog.FromEcho(c).Error(xlogger.Msgf("hello %s", "world!"), "first", first, "second", second)
			xlog.FromEcho(c).Error("oh snap! got error", "error", err)
	*/
	Error(msg string, fields ...any)

	/*
		Fields is a helper function to use a map or slice to set fields using type assertion.
		[]any must alternate string keys and arbitrary values, and extraneous ones are ignored. i.e:

		With Request HTTP Context:

			var (
				first  = "first value"
				second = "second value"
			)
			xlog.FromEcho(c).Fatal("hello", "first", first, "second", second)
			xlog.FromEcho(c).Fatal(xlogger.Msgf("hello %s", "world!"), "first", first, "second", second)
			xlog.FromEcho(c).Fatal("oh snap! got error", "error", err)
	*/
	Fatal(msg string, fields ...any)

	/*
		Fields is a helper function to use a map or slice to set fields using type assertion.
		[]any must alternate string keys and arbitrary values, and extraneous ones are ignored. i.e:

		With Request HTTP Context:

			var (
				first  = "first value"
				second = "second value"
			)
			xlog.FromEcho(c).Panic("hello", "first", first, "second", second)
			xlog.FromEcho(c).Panic(xlogger.Msgf("hello %s", "world!"), "first", first, "second", second)
			xlog.FromEcho(c).Panic("oh snap! got error", "error", err)
	*/
	Panic(msg string, fields ...any)
}

func Msgf(msg string, args ...any) string {
	return fmt.Sprintf(msg, args...)
}

type ZeroLogger struct {
	log zerolog.Logger
}

func NewLogger(log zerolog.Logger) Logger {
	return &ZeroLogger{log}
}

// attachFields directly to the zerolog.Event object without creating a new map
func (zl *ZeroLogger) attachFields(event *zerolog.Event, fields []any) *zerolog.Event {
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				event = event.Interface(key, fields[i+1])
			}
		}
	}
	return event
}

func (zl *ZeroLogger) Trace(msg string, fields ...any) {
	event := zl.log.Trace()
	event = zl.attachFields(event, fields)
	event.Msg(msg)
}

func (zl *ZeroLogger) Debug(msg string, fields ...any) {
	event := zl.log.Debug()
	event = zl.attachFields(event, fields)
	event.Msg(msg)
}

func (zl *ZeroLogger) Info(msg string, fields ...any) {
	event := zl.log.Info()
	event = zl.attachFields(event, fields)
	event.Msg(msg)
}

func (zl *ZeroLogger) Warn(msg string, fields ...any) {
	event := zl.log.Warn()
	event = zl.attachFields(event, fields)
	event.Msg(msg)
}

func (zl *ZeroLogger) Error(msg string, fields ...any) {
	event := zl.log.Error()
	event = zl.attachFields(event, fields)
	event.Msg(msg)
}

func (zl *ZeroLogger) Fatal(msg string, fields ...any) {
	event := zl.log.Fatal()
	event = zl.attachFields(event, fields)
	event.Msg(msg)
}

func (zl *ZeroLogger) Panic(msg string, fields ...any) {
	event := zl.log.Panic()
	event = zl.attachFields(event, fields)
	event.Msg(msg)
}
