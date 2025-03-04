package xlog

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/guregu/null/v5"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/log"
)

// Hook is a [zerolog.Hook] that sends all logging records it receives to
// OpenTelemetry. See package documentation for how conversions are made.
type OtelHook struct {
	lvl               zerolog.Level
	logWriterDisabled null.Bool
	logEnabled        null.Bool
	logger            log.Logger
}

// NewHook returns a new [Hook] to be used as a [Zerolog.Hook].
// If [WithLoggerProvider] is not provided, the returned Hook will use the
// global LoggerProvider.
func NewOtelHook(name string, options ...Option) *OtelHook {
	cfg := newConfig(options)
	return &OtelHook{
		logger:            cfg.logger(name),
		lvl:               cfg.level,
		logEnabled:        cfg.logEnabled,
		logWriterDisabled: cfg.logWriterDisabled,
	}
}

// Run handles the passed record, and sends it to OpenTelemetry.
func (h OtelHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// skip when log disabled
	if h.logEnabled.Valid {
		if !h.logEnabled.Bool {
			// fmt.Println("skiped log disabled", level, msg)
			return
		}
	}

	// skip when log writer disabled
	if h.logWriterDisabled.Valid && h.logWriterDisabled.Bool {
		// fmt.Println("skiped log writer disabled", level, msg)
		return
	}

	// skip level writer
	if level < h.lvl {
		// fmt.Println("skiped level", level, msg)
		return
	}

	var (
		rec  log.Record
		ctx  = e.GetCtx()
		data = make(map[string]any)
		p    = fmt.Appendf(nil, "%s}", reflect.ValueOf(e).Elem().FieldByName("buf"))
		_    = json.Unmarshal(p, &data)
	)

	if ctx == nil {
		ctx = context.Background()
	}

	data["message"] = msg

	var (
		attrs     = mapToKeyValues(data)
		buff, _   = json.Marshal(data)
		_time     = gjson.GetBytes(buff, "time").String()
		_vtime, _ = time.Parse(time.RFC3339Nano, _time)
		_lvl      = convertLevel(level)
	)

	rec.SetSeverity(_lvl)
	rec.SetBody(log.BytesValue(buff))
	rec.SetSeverityText(_lvl.String())
	rec.SetTimestamp(_vtime)
	rec.AddAttributes(attrs...)

	h.logger.Emit(ctx, rec)
}

func convertLevel(level zerolog.Level) log.Severity {
	switch level {
	case zerolog.DebugLevel:
		return log.SeverityDebug
	case zerolog.InfoLevel:
		return log.SeverityInfo
	case zerolog.WarnLevel:
		return log.SeverityWarn
	case zerolog.ErrorLevel:
		return log.SeverityError
	case zerolog.PanicLevel:
		return log.SeverityFatal1
	case zerolog.FatalLevel:
		return log.SeverityFatal2
	default:
		return log.SeverityUndefined
	}
}
