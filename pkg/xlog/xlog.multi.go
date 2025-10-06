package xlog

import (
	"errors"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// # Multi Logger
type Entry struct {
	key      string
	rotation *lumberjack.Logger
	Options  LogOptions
}

// Note:
//   - SetLogFileOutput(io.Writer) this option is not used in this option.
func NewEntry(key string, rotation *lumberjack.Logger, opts ...LogOptionFn) Entry {
	e := Entry{
		key:      key,
		rotation: rotation,
		Options: LogOptions{
			LogFields: make(map[string]any),
		},
	}

	for _, fn := range opts {
		fn(&e.Options)
	}

	return e
}

type MultiLogger struct {
	entryLogger   map[string]zerolog.Logger
	entryRotation map[string]*lumberjack.Logger
}

func NewMultiLogging(entries ...Entry) MultiLogger {
	l := MultiLogger{
		entryLogger:   make(map[string]zerolog.Logger),
		entryRotation: make(map[string]*lumberjack.Logger),
	}

	disabledFn := func(validity bool, value bool) bool {
		if validity {
			return value
		}
		return true
	}

	for _, e := range entries {
		var (
			isLogConsoleDisabled = disabledFn(e.Options.LogConsoleDisable.Valid, e.Options.LogConsoleDisable.Bool)
			isLogFileDisabled    = disabledFn(e.Options.LogFileDisable.Valid, e.Options.LogFileDisable.Bool)
		)

		if isLogConsoleDisabled || isLogFileDisabled {
			continue
		}

		if _, ok := l.entryRotation[e.key]; !ok {

			var (
				opts = []LogOptionFn{
					// Log Hook
					SetLogHook(e.Options.Hook...),

					// Log Console
					SetLogConsoleFormat(e.Options.LogConsoleFormat),
					SetLogConsoleDisabled(isLogConsoleDisabled),
					SetLogConsoleLevel(e.Options.LogConsoleLevel),
					SetLogConsoleOutput(e.Options.LogConsoleOut),

					// Log File
					SetLogFileDisabled(isLogFileDisabled),
					SetLogFileLevel(e.Options.LogFileLevel),
					SetLogFileOutput(e.rotation),
				}
			)

			l.entryRotation[e.key] = e.rotation
			l.entryLogger[e.key] = NewZeroLog(opts...)
			continue
		}
	}

	return l
}

func (c *MultiLogger) LogBy(key string) zerolog.Logger {
	if logger, ok := c.entryLogger[key]; ok {
		return logger
	}

	return NoopZeroLogger
}

func (c *MultiLogger) RotateBy(key string) error {
	if cfg, ok := c.entryRotation[key]; ok {
		return cfg.Rotate()
	}

	return errors.New("log key not found")
}

func (c *MultiLogger) Rotate() error {
	var errs []error

	for _, cfg := range c.entryRotation {
		if err := cfg.Rotate(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (c *MultiLogger) Close() {
	clear(c.entryRotation)
	clear(c.entryLogger)
}
