package xlog

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	SingleLogger struct {
		Logger      zerolog.Logger
		LogRotation lumberjack.Logger
	}

	DebugLogger struct {
		SingleLogger
	}

	IOLogger struct {
		SingleLogger
	}

	TRXLogger struct {
		MultiLogger
	}
)

func NoopSingleLogger() SingleLogger {
	return SingleLogger{Logger: nopZeroLogger}
}
