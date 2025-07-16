package xlog

import (
	"context"
	"encoding/json"
	"fmt"

	"resty.dev/v3"
)

type RestyV3Logger struct {
	Log Logger
}

func NewRestyV3Logger(log Logger) *RestyV3Logger {
	return &RestyV3Logger{log}
}

func (l *RestyV3Logger) Errorf(format string, v ...any) {
	var (
		ctx     = context.Background()
		logText = fmt.Sprintf(format, v...)
		fields  = []any{
			"restyLog", logText,
		}
	)

	l.Log.Error(ctx, "resty api log", fields...)
}

func (l *RestyV3Logger) Warnf(format string, v ...any) {
	var (
		ctx  = context.Background()
		text = fmt.Sprintf(format, v...)
	)

	var (
		fields = []any{}
		log    = resty.DebugLog{}
	)

	if err := json.Unmarshal([]byte(text), &log); err == nil && log.Request != nil && log.Response != nil {
		fields = append(fields, "restyParseLog", log)
	} else {
		fields = append(fields, "restyLog", text)
	}

	l.Log.Warn(ctx, "resty api log", fields...)
}

func (l *RestyV3Logger) Debugf(format string, v ...any) {
	var (
		ctx  = context.Background()
		text = fmt.Sprintf(format, v...)
	)

	var (
		fields = []any{}
		log    = resty.DebugLog{}
	)

	if err := json.Unmarshal([]byte(text), &log); err == nil && log.Request != nil && log.Response != nil {
		fields = append(fields, "restyParseLog", log)
	} else {
		fields = append(fields, "restyLog", text)
	}

	l.Log.Debug(ctx, "resty api log", fields...)
}
