package dependency

import (
	"path"
	"strings"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// DONT USE THIS FOR LOGGING
	TRXLogger *xlog.TRXLogger
)

func ProvideTRXLogger(c config.Cfg) *xlog.TRXLogger {
	var (
		l       xlog.TRXLogger
		entries = make([]xlog.Entry, len(c.Log.TRX.ClientKey))
	)

	for i, key := range c.Log.TRX.ClientKey {
		var (
			basePath = strings.ReplaceAll(
				c.Log.TRX.Rotation.BasePath,
				"{CLIENT}",
				strings.ToLower(
					strings.TrimSpace(key),
				),
			)

			filename = strings.ReplaceAll(
				c.Log.TRX.Rotation.Filename,
				"{CLIENT}",
				strings.ToLower(
					strings.TrimSpace(key),
				),
			)
		)

		entries[i] = xlog.NewEntryWithOptions(
			key,
			&lumberjack.Logger{
				Filename:   path.Join(basePath, filename), // where you need store to store log and what a log name
				MaxBackups: c.Log.TRX.Rotation.MaxBackup,  // how much backup files
				MaxSize:    c.Log.TRX.Rotation.MaxSize,    // how much maximum megabytes
				MaxAge:     c.Log.TRX.Rotation.MaxAge,     // how much maximum days, default is 0 that means not deleted old logs
				LocalTime:  c.Log.TRX.Rotation.LocalTime,  // default UTC or false
				Compress:   c.Log.TRX.Rotation.Compress,
			},
			xlog.EnableConsoleLogging(),
			xlog.SetLevelLog(c.Log.TRX.Level),
			xlog.SetField("appName", c.App.Name),
			xlog.SetField("appEnv", c.App.Env),
			xlog.SetField("appDomain", c.App.Domain),
			xlog.SetField("appLogKind", "trx-logger"),
			xlog.SetField("appLogTrx", key),
		)
	}

	l.MultiLogger = xlog.NewMultiLogging(entries...)
	TRXLogger = &l

	return &l
}
