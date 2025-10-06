package xlog

import (
	"io"
	"os"
	"time"

	"github.com/guregu/null/v5"
	"github.com/rs/zerolog"
)

type (
	LogOptions struct {
		// Log Hook
		Hook []zerolog.Hook

		// Log Fields
		LogFields map[string]any

		// Console
		LogConsoleFormat  string
		LogConsoleDisable null.Bool
		LogConsoleLevel   int
		LogConsoleOut     io.Writer

		// File
		LogFileDisable null.Bool
		LogFileLevel   int
		LogFileOut     io.Writer
	}

	LogOptionFn func(opt *LogOptions)
)

func SetField(key string, val any) LogOptionFn {
	return func(opt *LogOptions) {
		if _, ok := opt.LogFields[key]; ok {
			return
		}

		opt.LogFields[key] = val
	}
}

func SetLogHook(v ...zerolog.Hook) LogOptionFn {
	return func(opt *LogOptions) {
		opt.Hook = v
	}
}

func SetLogConsoleFormat(v string) LogOptionFn {
	return func(opt *LogOptions) {
		switch v {
		case "json":
			opt.LogConsoleFormat = "json"
		case "console":
			opt.LogConsoleFormat = "console"
		default:
			opt.LogConsoleFormat = "console"
		}
	}
}

func SetLogConsoleDisabled(v bool) LogOptionFn {
	return func(opt *LogOptions) {
		opt.LogConsoleDisable = null.BoolFrom(v)
	}
}

func SetLogConsoleLevel(lvl int) LogOptionFn {
	return func(opt *LogOptions) {
		opt.LogConsoleLevel = lvl
	}
}

func SetLogConsoleOutput(out io.Writer) LogOptionFn {
	return func(opt *LogOptions) {
		opt.LogConsoleOut = out
	}
}

func SetLogFileDisabled(v bool) LogOptionFn {
	return func(opt *LogOptions) {
		opt.LogFileDisable = null.BoolFrom(v)
	}
}

func SetLogFileLevel(lvl int) LogOptionFn {
	return func(opt *LogOptions) {
		opt.LogFileLevel = lvl
	}
}
func SetLogFileOutput(out io.Writer) LogOptionFn {
	return func(opt *LogOptions) {
		opt.LogFileOut = out
	}
}

func NewZeroLog(opts ...LogOptionFn) zerolog.Logger {
	var (
		mw  []io.Writer
		opt = LogOptions{
			LogFields:        make(map[string]any),
			LogConsoleOut:    os.Stderr,
			LogConsoleFormat: "console",
		}
	)

	for _, fn := range opts {
		fn(&opt)
	}

	disabledFn := func(validity bool, value bool) bool {
		if validity {
			return value
		}
		return true
	}

	var (
		isLogConsoleDisabled = disabledFn(opt.LogConsoleDisable.Valid, opt.LogConsoleDisable.Bool)
		isLogFileDisabled    = disabledFn(opt.LogFileDisable.Valid, opt.LogFileDisable.Bool)
	)

	if isLogConsoleDisabled && isLogFileDisabled {
		return zerolog.Nop()
	}

	if !isLogConsoleDisabled {
		var w io.Writer
		switch opt.LogConsoleFormat {
		case "json":
			w = opt.LogConsoleOut
		case "console":
			w = zerolog.ConsoleWriter{Out: opt.LogConsoleOut}
		}

		consoleLog := zerolog.FilteredLevelWriter{
			Writer: zerolog.LevelWriterAdapter{Writer: w},
			Level:  zerolog.Level(opt.LogConsoleLevel),
		}

		mw = append(mw, &consoleLog)
	}

	if !isLogFileDisabled {
		fileLog := zerolog.FilteredLevelWriter{
			Writer: zerolog.LevelWriterAdapter{Writer: opt.LogFileOut},
			Level:  zerolog.Level(opt.LogFileLevel),
		}

		// Set Log File
		mw = append(mw, &fileLog)
	}

	ctx := zerolog.New(zerolog.MultiLevelWriter(mw...)).With()
	if len(opt.LogFields) > 0 {
		ctx = AnyFieldsToContext(ctx, opt.LogFields)
	}

	// Set Default into time format with nano
	zerolog.TimeFieldFormat = time.RFC3339Nano

	return ctx.Timestamp().Logger().Hook(opt.Hook...)
}
