package dependency

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"go.uber.org/fx"

	"github.com/rs/zerolog"
	otelog "go.opentelemetry.io/otel/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	isFileLogDisabled bool
	debugLogger       *xlog.DebugLogger
)

func RotateLog() {
	if debugLogger != nil && !isFileLogDisabled {
		if err := debugLogger.LogRotation.Rotate(); err != nil {
			log.Printf("try to rotate 'debug:logger', got err: %+v\n", err)
		}
	}
}

type DebugLoggerParamFx struct {
	fx.In
	Cfg            config.Cfg
	Server         config.Server
	LoggerProvider otelog.LoggerProvider `optional:"true"`
}

func ProvideDebugLogger(p DebugLoggerParamFx) *xlog.DebugLogger {
	var (
		filename          = path.Join(p.Cfg.Log.BasePath, p.Cfg.Log.File.Name)
		rotation          = p.Cfg.Log.File.Rotation
		isFileLogDisabled = p.Cfg.Log.File.Disabled
		debugLog          = xlog.DebugLogger{
			SingleLogger: xlog.SingleLogger{
				LogRotation: &lumberjack.Logger{
					Filename:   filename,
					MaxBackups: rotation.MaxBackup,
					MaxSize:    rotation.MaxSize,
					MaxAge:     rotation.MaxAge,
					LocalTime:  rotation.LocalTime,
					Compress:   rotation.Compress,
				},
			},
		}
	)

	debugLog.Logger = xlog.NewZeroLog(
		// otel options
		xlog.SetLogHook(
			xlog.NewOtelHook(
				fmt.Sprintf("%s/%s", p.Cfg.App.Project, p.Server.Name),
				xlog.WithLogEnabled(p.Cfg.Otel.Logs),
				xlog.WithListIgnoreKeys(p.Cfg.Otel.Options.IgnoreLogKeys),
				xlog.WithLogWriterDisabled(!p.Cfg.Otel.Logs),
				xlog.WithLevel(zerolog.Level(p.Cfg.Log.Level)),
				xlog.WithLoggerProvider(p.LoggerProvider),
				xlog.WithVersion("1.0.0"),
			),
		),

		// console log
		xlog.SetLogConsoleFormat(p.Cfg.Log.ConsoleFormat),
		xlog.SetLogConsoleDisabled(false),
		xlog.SetLogConsoleLevel(p.Cfg.Log.Level),
		xlog.SetLogConsoleOutput(os.Stderr),

		// file log
		xlog.SetLogFileDisabled(isFileLogDisabled),
		xlog.SetLogFileLevel(p.Cfg.Log.Level),
		xlog.SetLogFileOutput(debugLog.LogRotation),

		// options fields
		xlog.SetField("appName", fmt.Sprintf("%s/%s", p.Cfg.App.Project, p.Server.Name)),
		xlog.SetField("appEnv", p.Cfg.App.Env),
		xlog.SetField("appServer", p.Server.Name),
		xlog.SetField("appLog", "debug:logger"),
		xlog.SetField("appPid", os.Getpid()),
	)

	debugLogger = &debugLog

	return &debugLog
}
