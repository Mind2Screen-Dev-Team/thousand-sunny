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
		xfilename = path.Join(cfg.Rotation.BasePath, cfg.Rotation.Filename)
		l         = xlog.DebugLogger{
			SingleLogger: xlog.SingleLogger{
				LogRotation: lumberjack.Logger{
					Filename:   xfilename,              // where you need store to store log and what a log name
					MaxBackups: cfg.Rotation.MaxBackup, // how much backup files
					MaxSize:    cfg.Rotation.MaxSize,   // how much maximum megabytes
					MaxAge:     cfg.Rotation.MaxAge,    // how much maximum days, default is 0 that means not deleted old logs
					LocalTime:  cfg.Rotation.LocalTime, // default UTC -> false
					Compress:   cfg.Rotation.Compress,  // default Un-Compressed -> false
				},
			},
		}
	)

	l.Logger = xlog.NewZeroLog(
		&l.LogRotation,

		// # options
		xlog.EnableConsoleLogging(),
		xlog.SetConsoleOutput(os.Stderr),
		xlog.SetLevelLog(cfg.Level),

		// # options fields
		xlog.SetField("appName", c.App.Name),
		xlog.SetField("appEnv", c.App.Env),
		xlog.SetField("appLog", "debug:logger"),
	)
	_DebugLogger = &l

	return &l
}

func ProvideIoLogger(c config.Cfg) *xlog.IOLogger {
	var (
		cfg       = c.Log["io"]
		xfilename = path.Join(cfg.Rotation.BasePath, cfg.Rotation.Filename)
		l         = xlog.IOLogger{
			SingleLogger: xlog.SingleLogger{
				LogRotation: lumberjack.Logger{
					Filename:   xfilename,              // where you need store to store log and what a log name
					MaxBackups: cfg.Rotation.MaxBackup, // how much backup files
					MaxSize:    cfg.Rotation.MaxSize,   // how much maximum megabytes
					MaxAge:     cfg.Rotation.MaxAge,    // how much maximum days, default is 0 that means not deleted old logs
					LocalTime:  cfg.Rotation.LocalTime, // default UTC or false
					Compress:   cfg.Rotation.Compress,  // default Un-Compressed or false
				},
			},
		}
	)

	l.Logger = xlog.NewZeroLog(
		&l.LogRotation,

		// # options
		xlog.EnableConsoleLogging(),
		xlog.SetConsoleOutput(os.Stderr),
		xlog.SetLevelLog(cfg.Level),

		// # fields
		xlog.SetField("appName", c.App.Name),
		xlog.SetField("appEnv", c.App.Env),
		xlog.SetField("appLog", "io:logger"),
	)
	_IOLogger = &l

	return &l
}

func ProvideTrxLogger(c config.Cfg) *xlog.TrxLogger {
	var (
		cfg     = c.Log["trx"]
		l       xlog.TrxLogger
		entries = make([]xlog.Entry, len(cfg.ClientKey))
	)

	for i, key := range cfg.ClientKey {
		var (
			basePath = strings.ReplaceAll(
				cfg.Rotation.BasePath,
				"{CLIENT}",
				strings.ToLower(
					strings.TrimSpace(key),
				),
			)

			filename = strings.ReplaceAll(
				cfg.Rotation.Filename,
				"{CLIENT}",
				strings.ToLower(
					strings.TrimSpace(key),
				),
			)

			xfilename = path.Join(basePath, filename)
		)

		entries[i] = xlog.NewEntryWithOptions(
			key,
			&lumberjack.Logger{
				Filename:   xfilename,              // where you need store to store log and what a log name
				MaxBackups: cfg.Rotation.MaxBackup, // how much backup files
				MaxSize:    cfg.Rotation.MaxSize,   // how much maximum megabytes
				MaxAge:     cfg.Rotation.MaxAge,    // how much maximum days, default is 0 that means not deleted old logs
				LocalTime:  cfg.Rotation.LocalTime, // default UTC or false
				Compress:   cfg.Rotation.Compress,
			},

			// # options
			xlog.EnableConsoleLogging(),
			xlog.SetConsoleOutput(os.Stderr),
			xlog.SetLevelLog(cfg.Level),

			// # fields
			xlog.SetField("appName", c.App.Name),
			xlog.SetField("appEnv", c.App.Env),
			xlog.SetField("appLog", fmt.Sprintf("trx:logger:%s", key)),
		)
	}

	l.MultiLogger = xlog.NewMultiLogging(entries...)
	_TrxLogger = &l

	return &l
}
