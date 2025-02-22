package xlog

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type (
	LogOptions struct {
		// Log Fields
		LogFields map[string]any

		// Console
		LogConsoleDisable bool
		LogConsoleLevel   int
		LogConsoleOut     io.Writer

		// File
		LogFileDisable bool
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

func SetLogConsoleDisabled(v bool) LogOptionFn {
	return func(opt *LogOptions) {
		opt.LogConsoleDisable = v
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
		opt.LogFileDisable = v
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
			LogFields:     make(map[string]any),
			LogConsoleOut: os.Stderr,
		}
	)

	for _, fn := range opts {
		fn(&opt)
	}

	if opt.LogConsoleDisable && opt.LogFileDisable {
		return zerolog.Nop()
	}

	if !opt.LogConsoleDisable {
		consoleLog := zerolog.FilteredLevelWriter{
			Writer: zerolog.LevelWriterAdapter{Writer: zerolog.ConsoleWriter{Out: opt.LogConsoleOut}},
			Level:  zerolog.Level(opt.LogConsoleLevel),
		}

		mw = append(mw, &consoleLog)
	}

	if !opt.LogFileDisable {
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

	return ctx.Timestamp().Logger()
}
