package dependency

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	_DebugLogger *xlog.DebugLogger
	_IOLogger    *xlog.IOLogger
	_TrxLogger   *xlog.TrxLogger
)

func RotateLog() {
	if _DebugLogger != nil {
		if err := _DebugLogger.LogRotation.Rotate(); err != nil {
			log.Printf("try to rotate 'debug:logger', got err: %+v\n", err)
		}
	}

	if _IOLogger != nil {
		if err := _IOLogger.LogRotation.Rotate(); err != nil {
			log.Printf("try to rotate 'io:logger', got err: %+v\n", err)
		}
	}

	if _TrxLogger != nil {
		if err := _TrxLogger.Rotate(); err != nil {
			log.Printf("try to rotate 'trx:logger', got err: %+v\n", err)
		}

		// clear memory trx logger
		_TrxLogger.Close()
	}
}

func ProvideDebugLogger(c config.Cfg) *xlog.DebugLogger {
	var (
		cfg       = c.Log["debug"]
		xfilename = path.Join(cfg.File.Rotation.BasePath, cfg.File.Rotation.Filename)
		debugLog  = xlog.DebugLogger{
			SingleLogger: xlog.SingleLogger{
				LogRotation: lumberjack.Logger{
					Filename:   xfilename,                   // where you need store to store log and what a log name
					MaxBackups: cfg.File.Rotation.MaxBackup, // how much backup files
					MaxSize:    cfg.File.Rotation.MaxSize,   // how much maximum megabytes
					MaxAge:     cfg.File.Rotation.MaxAge,    // how much maximum days, default is 0 that means not deleted old logs
					LocalTime:  cfg.File.Rotation.LocalTime, // default UTC -> false
					Compress:   cfg.File.Rotation.Compress,  // default Un-Compressed -> false
				},
			},
		}
	)

	debugLog.Logger = xlog.NewZeroLog(
		// # options
		xlog.SetLogConsoleDisabled(cfg.Console.Disabled),
		xlog.SetLogConsoleLevel(cfg.Console.Level),
		xlog.SetLogConsoleOutput(os.Stderr),

		xlog.SetLogFileDisabled(cfg.File.Disabled),
		xlog.SetLogFileLevel(cfg.File.Level),
		xlog.SetLogFileOutput(&debugLog.LogRotation),

		// # Options Fields
		xlog.SetField("appName", c.App.Name),
		xlog.SetField("appEnv", c.App.Env),
		xlog.SetField("appLog", "debug:logger"),
	)

	_DebugLogger = &debugLog

	return &debugLog
}

func ProvideIoLogger(c config.Cfg) *xlog.IOLogger {
	var (
		cfg       = c.Log["io"]
		xfilename = path.Join(cfg.File.Rotation.BasePath, cfg.File.Rotation.Filename)
		ioLog     = xlog.IOLogger{
			SingleLogger: xlog.SingleLogger{
				LogRotation: lumberjack.Logger{
					Filename:   xfilename,                   // where you need store to store log and what a log name
					MaxBackups: cfg.File.Rotation.MaxBackup, // how much backup files
					MaxSize:    cfg.File.Rotation.MaxSize,   // how much maximum megabytes
					MaxAge:     cfg.File.Rotation.MaxAge,    // how much maximum days, default is 0 that means not deleted old logs
					LocalTime:  cfg.File.Rotation.LocalTime, // default UTC or false
					Compress:   cfg.File.Rotation.Compress,  // default Un-Compressed or false
				},
			},
		}
	)

	ioLog.Logger = xlog.NewZeroLog(
		// # options
		xlog.SetLogConsoleDisabled(cfg.Console.Disabled),
		xlog.SetLogConsoleLevel(cfg.Console.Level),
		xlog.SetLogConsoleOutput(os.Stderr),

		xlog.SetLogFileDisabled(cfg.File.Disabled),
		xlog.SetLogFileLevel(cfg.File.Level),
		xlog.SetLogFileOutput(&ioLog.LogRotation),

		// # fields
		xlog.SetField("appName", c.App.Name),
		xlog.SetField("appEnv", c.App.Env),
		xlog.SetField("appLog", "io:logger"),
	)
	_IOLogger = &ioLog

	return &ioLog
}

func ProvideTrxLogger(c config.Cfg) *xlog.TrxLogger {
	var (
		// TODO: you can setup this from db or configs
		clientKey = []string{}
	)

	var (
		cfg     = c.Log["trx"]
		l       xlog.TrxLogger
		entries = make([]xlog.Entry, 0)
	)

	for i, key := range clientKey {
		var (
			basePath = strings.ReplaceAll(
				cfg.File.Rotation.BasePath,
				"{CLIENT}",
				strings.ToLower(
					strings.TrimSpace(key),
				),
			)

			filename = strings.ReplaceAll(
				cfg.File.Rotation.Filename,
				"{CLIENT}",
				strings.ToLower(
					strings.TrimSpace(key),
				),
			)

			xfilename = path.Join(basePath, filename)

			logFileRotation = lumberjack.Logger{
				Filename:   xfilename,                   // where you need store to store log and what a log name
				MaxBackups: cfg.File.Rotation.MaxBackup, // how much backup files
				MaxSize:    cfg.File.Rotation.MaxSize,   // how much maximum megabytes
				MaxAge:     cfg.File.Rotation.MaxAge,    // how much maximum days, default is 0 that means not deleted old logs
				LocalTime:  cfg.File.Rotation.LocalTime, // default UTC or false
				Compress:   cfg.File.Rotation.Compress,  // default Un-Compressed -> false
			}
		)

		entries[i] = xlog.NewEntry(
			key,
			&logFileRotation,

			// # options
			xlog.SetLogConsoleDisabled(cfg.Console.Disabled),
			xlog.SetLogConsoleLevel(cfg.Console.Level),
			xlog.SetLogConsoleOutput(os.Stderr),

			xlog.SetLogFileDisabled(cfg.File.Disabled),
			xlog.SetLogFileLevel(cfg.File.Level),
			xlog.SetLogFileOutput(&logFileRotation),

			// # fields
			xlog.SetField("appName", c.App.Name),
			xlog.SetField("appEnv", c.App.Env),
			xlog.SetField("appLog", fmt.Sprintf("trx:logger:%s", key)),
		)
	}

	l.MultiLogger = xlog.NewMultiLogging(cfg.File.Level, entries...)
	_TrxLogger = &l

	return &l
}
