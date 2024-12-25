package dependency

import (
	"path"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
)

var (
	// DONT USE THIS FOR LOGGING
	IOLogger *xlog.IOLogger
)

func ProvideIOLogger(c config.Cfg) *xlog.IOLogger {
	l := xlog.IOLogger{
		SingleLogger: xlog.SingleLogger{
			LogRotation: lumberjack.Logger{
				Filename:   path.Join(c.Log.IO.Rotation.BasePath, c.Log.IO.Rotation.Filename), // where you need store to store log and what a log name
				MaxBackups: c.Log.IO.Rotation.MaxBackup,                                       // how much backup files
				MaxSize:    c.Log.IO.Rotation.MaxSize,                                         // how much maximum megabytes
				MaxAge:     c.Log.IO.Rotation.MaxAge,                                          // how much maximum days, default is 0 that means not deleted old logs
				LocalTime:  c.Log.IO.Rotation.LocalTime,                                       // default UTC or false
				Compress:   c.Log.IO.Rotation.Compress,                                        // default Un-Compressed or false
			},
		},
	}

	l.Logger = xlog.NewZeroLog(
		&l.LogRotation,
		xlog.EnableConsoleLogging(),
		xlog.SetLevelLog(c.Log.IO.Level),
		xlog.SetField("appName", c.App.Name),
		xlog.SetField("appEnv", c.App.Env),
		xlog.SetField("appDomain", c.App.Domain),
		xlog.SetField("appLogKind", "io-logger"),
	)

	IOLogger = &l

	return &l
}
