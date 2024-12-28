package xlog

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogOptions struct {
	// Fields
	Fields map[string]any

	// Show log at console
	EnabledConsole bool

	// Level creates a child logger with the minimum accepted level set to level.
	Level int

	// Console Out
	Out io.Writer
}

type LogOptionFn func(opt *LogOptions)

func SetField(key string, val any) LogOptionFn {
	return func(opt *LogOptions) {
		if _, ok := opt.Fields[key]; ok {
			return
		}

		opt.Fields[key] = val
	}
}

func EnableConsoleLogging() LogOptionFn {
	return func(opt *LogOptions) {
		opt.EnabledConsole = true
	}
}

func SetConsoleOutput(out io.Writer) LogOptionFn {
	return func(opt *LogOptions) {
		opt.Out = out
	}
}

func SetLevelLog(lvl int) LogOptionFn {
	return func(opt *LogOptions) {
		opt.Level = lvl
	}
}

func NewZeroLog(rotation *lumberjack.Logger, opts ...LogOptionFn) zerolog.Logger {
	var (
		mw  []io.Writer
		lvl zerolog.Level
		opt = LogOptions{
			Fields: make(map[string]any),
			Out:    os.Stderr,
		}
	)

	for _, fn := range opts {
		fn(&opt)
	}

	if opt.EnabledConsole {
		mw = append(mw, zerolog.ConsoleWriter{Out: opt.Out})
	}

	// Set Log File
	mw = append(mw, rotation)

	switch opt.Level {
	default:
		lvl = zerolog.TraceLevel
	case 0:
		lvl = zerolog.DebugLevel
	case 1:
		lvl = zerolog.InfoLevel
	case 2:
		lvl = zerolog.WarnLevel
	case 3:
		lvl = zerolog.ErrorLevel
	case 4:
		lvl = zerolog.FatalLevel
	case 5:
		lvl = zerolog.PanicLevel
	}

	ctx := zerolog.New(io.MultiWriter(mw...)).Level(lvl).With()
	if len(opt.Fields) > 0 {
		ctx = ctx.Fields(opt.Fields)
	}

	// Set Default into time format with nano
	zerolog.TimeFieldFormat = time.RFC3339Nano

	return ctx.Timestamp().Logger()
}

// # Multi Logger

type Entry struct {
	key string
	cfg *lumberjack.Logger

	// options
	lvl            int
	enabledConsole bool
}

func NewEntry(key string, cfg *lumberjack.Logger) Entry {
	return Entry{
		key:            key,
		cfg:            cfg,
		lvl:            int(zerolog.TraceLevel),
		enabledConsole: false,
	}
}

func NewEntryWithOptions(key string, cfg *lumberjack.Logger, opts ...LogOptionFn) Entry {
	var (
		opt = LogOptions{}
		e   = Entry{key: key, cfg: cfg}
	)

	for _, fn := range opts {
		fn(&opt)
	}

	e.lvl = opt.Level
	e.enabledConsole = opt.EnabledConsole

	return e
}

type MultiLogger struct {
	entryLogger map[string]zerolog.Logger
	entryCfg    map[string]*lumberjack.Logger
}

func NewMultiLogging(entries ...Entry) MultiLogger {
	l := MultiLogger{
		entryLogger: make(map[string]zerolog.Logger),
		entryCfg:    make(map[string]*lumberjack.Logger),
	}

	for _, e := range entries {
		if _, ok := l.entryCfg[e.key]; !ok {

			opts := []LogOptionFn{
				SetLevelLog(e.lvl),
			}

			if e.enabledConsole {
				opts = append(opts, EnableConsoleLogging())
			}

			l.entryCfg[e.key] = e.cfg
			l.entryLogger[e.key] = NewZeroLog(e.cfg, opts...)
			continue
		}
	}

	return l
}

func (c *MultiLogger) LogBy(key string) zerolog.Logger {
	if logger, ok := c.entryLogger[key]; ok {
		return logger
	}

	return nopZeroLogger
}

func (c *MultiLogger) RotateBy(key string) error {
	if cfg, ok := c.entryCfg[key]; ok {
		return cfg.Rotate()
	}

	return errors.New("log key not found")
}

func (c *MultiLogger) Rotate() error {
	var errs []error

	for _, cfg := range c.entryCfg {
		if err := cfg.Rotate(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (c *MultiLogger) Close() {
	clear(c.entryCfg)
	clear(c.entryLogger)
}
