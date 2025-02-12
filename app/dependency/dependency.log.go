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

type Server string

func (s Server) String() string {
	return string(s)
}

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

func ProvideDebugLogger(c config.Cfg, s Server) *xlog.DebugLogger {

	basePath := strings.ReplaceAll(c.Log.BasePath, "{server_name}", s.String())
	basePath = strings.ReplaceAll(basePath, "{log_type}", "debug")

	var (
		cfg       = c.Log.LogType["debug"]
		xfilename = path.Join(basePath, cfg.File.Rotation.Filename)
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
		xlog.SetField("appServer", s.String()),
		xlog.SetField("appLog", "debug:logger"),
	)

	_DebugLogger = &debugLog

	return &debugLog
}

func ProvideIoLogger(c config.Cfg, s Server) *xlog.IOLogger {

	basePath := strings.ReplaceAll(c.Log.BasePath, "{server_name}", s.String())
	basePath = strings.ReplaceAll(basePath, "{log_type}", "io")

	var (
		cfg       = c.Log.LogType["io"]
		xfilename = path.Join(basePath, cfg.File.Rotation.Filename)
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
		xlog.SetField("appServer", s.String()),
		xlog.SetField("appLog", "io:logger"),
	)
	_IOLogger = &ioLog

	return &ioLog
}

func ProvideTrxLogger(c config.Cfg, s Server) *xlog.TrxLogger {
	basePath := strings.ReplaceAll(c.Log.BasePath, "{server_name}", s.String())
	basePath = strings.ReplaceAll(basePath, "{log_type}", "trx")
	basePath = strings.Join([]string{basePath, "{client}"}, "/")

	var (
		l       = xlog.TrxLogger{}
		cfg     = c.Log.LogType["trx"]
		entries = make([]xlog.Entry, len(c.Log.Client))
	)

	for i, key := range c.Log.Client {
		var (
			k               = strings.Split(strings.ToLower(strings.TrimSpace(key)), ":")
			basePath        = strings.ReplaceAll(basePath, "{client}", strings.Join(k, "/"))
			xfilename       = path.Join(basePath, cfg.File.Rotation.Filename)
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
			xlog.SetField("appServer", s.String()),
			xlog.SetField("appLog", fmt.Sprintf("trx:logger:%s", key)),
		)
	}

	l.MultiLogger = xlog.NewMultiLogging(cfg.File.Level, entries...)
	_TrxLogger = &l

	return &l
}
