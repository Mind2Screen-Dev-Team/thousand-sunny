package xlog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	provider   log.LoggerProvider
	version    string
	schemaURL  string
	serverName string
	serverAddr string
}

func newConfig(options []Option) config {
	var c config
	for _, opt := range options {
		c = opt.apply(c)
	}

	if c.provider == nil {
		c.provider = global.GetLoggerProvider()
	}
	return c
}

func (c config) logger(name string) log.Logger {
	var (
		opts []log.LoggerOption
		attr []attribute.KeyValue
	)

	if c.serverName != "" {
		attr = append(attr, attribute.String("server.name", c.serverName))
	}

	if c.serverName != "" {
		attr = append(attr, attribute.String("server.addr", c.serverAddr))
	}

	if len(attr) > 0 {
		opts = append(opts, log.WithInstrumentationAttributes(attr...))
	}

	if c.version != "" {
		opts = append(opts, log.WithInstrumentationVersion(c.version))
	}

	if c.schemaURL != "" {
		opts = append(opts, log.WithSchemaURL(c.schemaURL))
	}

	return c.provider.Logger(name, opts...)
}

// Option configures a Hook.
type Option interface {
	apply(config) config
}
type optFunc func(config) config

func (f optFunc) apply(c config) config { return f(c) }

// WithVersion returns an [Option] that configures the version of the
// [log.Logger] used by a [OtelLogger]. The version should be the version of the
// package that is being logged.
func WithVersion(version string) Option {
	return optFunc(func(c config) config {
		c.version = version
		return c
	})
}

// WithSchemaURL returns an [Option] that configures the semantic convention
// schema URL of the [log.Logger] used by a [OtelLogger]. The schemaURL should be
// the schema URL for the semantic conventions used in log records.
func WithSchemaURL(schemaURL string) Option {
	return optFunc(func(c config) config {
		c.schemaURL = schemaURL
		return c
	})
}

func WithServerName(serverName string) Option {
	return optFunc(func(c config) config {
		c.serverName = serverName
		return c
	})
}

func WithServerAddress(serverAddres string) Option {
	return optFunc(func(c config) config {
		c.serverAddr = serverAddres
		return c
	})
}

// WithLoggerProvider returns an [Option] that configures [log.LoggerProvider]
//
// By default if this Option is not provided, the Logger will use the global LoggerProvider.
func WithLoggerProvider(provider log.LoggerProvider) Option {
	return optFunc(func(c config) config {
		c.provider = provider
		return c
	})
}

var _ io.WriteCloser = (*OtelLogger)(nil)

type OtelLogger struct {
	logger log.Logger
}

func NewOtelLogger(name string, options ...Option) *OtelLogger {
	cfg := newConfig(options)
	return &OtelLogger{
		logger: cfg.logger(name),
	}
}

func (l *OtelLogger) Write(p []byte) (n int, err error) {
	var (
		rec               log.Record
		ctx               = context.Background()
		lvl               = log.SeverityInfo
		data              = make(map[string]any)
		_                 = json.Unmarshal(p, &data)
		attrs, span, _err = mapToKeyValues(data)

		_lvl      = strings.ToUpper(gjson.GetBytes(p, "level").String())
		_msg      = gjson.GetBytes(p, "message").String()
		_time     = gjson.GetBytes(p, "time").String()
		_vtime, _ = time.Parse(time.RFC3339Nano, _time)
	)

	if _err != nil {
		ctx = trace.ContextWithSpanContext(ctx, span)
	}

	switch {
	case _lvl == log.SeverityTrace.String():
		lvl = log.SeverityTrace
	case _lvl == log.SeverityDebug.String():
		lvl = log.SeverityDebug
	case _lvl == log.SeverityInfo.String():
		lvl = log.SeverityInfo
	case _lvl == log.SeverityWarn.String():
		lvl = log.SeverityWarn
	case _lvl == log.SeverityError.String():
		lvl = log.SeverityError
	case _lvl == log.SeverityFatal.String():
		lvl = log.SeverityFatal
	}

	rec.SetSeverity(lvl)
	rec.SetBody(log.StringValue(_msg))
	rec.SetSeverityText(lvl.String())
	rec.SetTimestamp(_vtime)
	rec.AddAttributes(attrs...)

	l.logger.Emit(ctx, rec)

	clear(data)

	return len(p), nil
}

// Close implements io.Closer, and closes the current logfile.
func (l *OtelLogger) Close() error {
	return nil
}

// Converts map[string]any to []log.KeyValue
func mapToKeyValues(input map[string]any) (attrs []log.KeyValue, span trace.SpanContext, err error) {
	result := make([]log.KeyValue, 0, len(input))

	var (
		_spanId, _      = input["otel_span_id"].(string)
		_traceId, _     = input["otel_trace_id"].(string)
		_traceFlags, _  = input["otel_trace_flags"].(byte)
		_traceState, _  = input["otel_trace_state"].(string)
		_traceRemote, _ = input["otel_trace_remote"].(bool)
	)

	for key, value := range input {
		result = append(result, log.KeyValue{
			Key:   key,
			Value: convertValue(value),
		})
	}

	var errs []error
	spanId, err := trace.SpanIDFromHex(_spanId)
	if err != nil {
		errs = append(errs, err)
	}

	traceId, err := trace.TraceIDFromHex(_traceId)
	if err != nil {
		errs = append(errs, err)
	}

	traceState, _ := trace.ParseTraceState(_traceState)

	if len(errs) > 0 {
		err = errors.Join(errs...)
	}

	return result, trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceId,
		SpanID:     spanId,
		TraceState: traceState,
		TraceFlags: trace.TraceFlags(_traceFlags),
		Remote:     _traceRemote,
	}), err
}

// convertValue converts various types to log.Value.
func convertValue(v any) log.Value {
	// Handling the most common types without reflect is a small perf win.
	switch val := v.(type) {
	case bool:
		return log.BoolValue(val)
	case string:
		return log.StringValue(val)
	case int:
		return log.Int64Value(int64(val))
	case int8:
		return log.Int64Value(int64(val))
	case int16:
		return log.Int64Value(int64(val))
	case int32:
		return log.Int64Value(int64(val))
	case int64:
		return log.Int64Value(val)
	case uint:
		return convertUintValue(uint64(val))
	case uint8:
		return log.Int64Value(int64(val))
	case uint16:
		return log.Int64Value(int64(val))
	case uint32:
		return log.Int64Value(int64(val))
	case uint64:
		return convertUintValue(val)
	case uintptr:
		return convertUintValue(uint64(val))
	case float32:
		return log.Float64Value(float64(val))
	case float64:
		return log.Float64Value(val)
	case time.Duration:
		return log.Int64Value(val.Nanoseconds())
	case complex64:
		r := log.Float64("r", real(complex128(val)))
		i := log.Float64("i", imag(complex128(val)))
		return log.MapValue(r, i)
	case complex128:
		r := log.Float64("r", real(val))
		i := log.Float64("i", imag(val))
		return log.MapValue(r, i)
	case time.Time:
		return log.Int64Value(val.UnixNano())
	case []byte:
		return log.BytesValue(val)
	case error:
		return log.StringValue(val.Error())
	}

	t := reflect.TypeOf(v)
	if t == nil {
		return log.Value{}
	}
	val := reflect.ValueOf(v)
	switch t.Kind() {
	case reflect.Struct:
		return log.StringValue(fmt.Sprintf("%+v", v))
	case reflect.Slice, reflect.Array:
		items := make([]log.Value, 0, val.Len())
		for i := 0; i < val.Len(); i++ {
			items = append(items, convertValue(val.Index(i).Interface()))
		}
		return log.SliceValue(items...)
	case reflect.Map:
		kvs := make([]log.KeyValue, 0, val.Len())
		for _, k := range val.MapKeys() {
			var key string
			switch k.Kind() {
			case reflect.String:
				key = k.String()
			default:
				key = fmt.Sprintf("%+v", k.Interface())
			}
			kvs = append(kvs, log.KeyValue{
				Key:   key,
				Value: convertValue(val.MapIndex(k).Interface()),
			})
		}
		return log.MapValue(kvs...)
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			return log.Value{}
		}
		return convertValue(val.Elem().Interface())
	}

	// Try to handle this as gracefully as possible.
	//
	// Don't panic here. it is preferable to have user's open issue
	// asking why their attributes have a "unhandled: " prefix than
	// say that their code is panicking.
	return log.StringValue(fmt.Sprintf("unhandled: (%s) %+v", t, v))
}

// convertUintValue converts a uint64 to a log.Value.
// If the value is too large to fit in an int64, it is converted to a string.
func convertUintValue(v uint64) log.Value {
	if v > math.MaxInt64 {
		return log.StringValue(strconv.FormatUint(v, 10))
	}
	return log.Int64Value(int64(v))
}
