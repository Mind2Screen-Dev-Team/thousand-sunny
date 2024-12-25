package dependency

import (
	"path"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
)

var (
	// DONT USE THIS FOR LOGGING
	DebugLogger *xlog.DebugLogger
)

func ProvideDebugLogger(c config.Cfg) *xlog.DebugLogger {
	l := xlog.DebugLogger{
		SingleLogger: xlog.SingleLogger{
			LogRotation: lumberjack.Logger{
				Filename:   path.Join(c.Log.Debug.Rotation.BasePath, c.Log.Debug.Rotation.Filename), // where you need store to store log and what a log name
				MaxBackups: c.Log.Debug.Rotation.MaxBackup,                                          // how much backup files
				MaxSize:    c.Log.Debug.Rotation.MaxSize,                                            // how much maximum megabytes
				MaxAge:     c.Log.Debug.Rotation.MaxAge,                                             // how much maximum days, default is 0 that means not deleted old logs
				LocalTime:  c.Log.Debug.Rotation.LocalTime,                                          // default UTC -> false
				Compress:   c.Log.Debug.Rotation.Compress,                                           // default Un-Compressed -> false
			},
		},
	}

	l.Logger = xlog.NewZeroLog(
		&l.LogRotation,
		xlog.EnableConsoleLogging(),
		xlog.SetLevelLog(c.Log.Debug.Level),
		xlog.SetField("appName", c.App.Name),
		xlog.SetField("appEnv", c.App.Env),
		xlog.SetField("appDomain", c.App.Domain),
		xlog.SetField("appLogKind", "debug-logger"),
	)

	DebugLogger = &l

	return &l
}
