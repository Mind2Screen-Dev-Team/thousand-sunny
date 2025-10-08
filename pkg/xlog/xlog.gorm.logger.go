package xlog

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	Log       Logger
	LogLevel  logger.LogLevel
	SlowQuery time.Duration
}

func NewGormLogger(log Logger, level logger.LogLevel, slow time.Duration) *GormLogger {
	return &GormLogger{
		Log:       log,
		LogLevel:  level,
		SlowQuery: slow,
	}
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	g.LogLevel = level
	return g
}

func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Info {
		g.Log.Debug(ctx, fmt.Sprintf(msg, data...))
	}
}

func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Warn {
		g.Log.Debug(ctx, fmt.Sprintf("WARN: "+msg, data...))
	}
}

func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Error {
		g.Log.Debug(ctx, fmt.Sprintf("ERROR: "+msg, data...))
	}
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.LogLevel == logger.Silent {
		return
	}

	sqlStr, rows := fc()
	elapsed := time.Since(begin)

	fields := make([]any, 0)
	if id, ok := ctx.Value(XLOG_REQ_TRACE_ID_CTX_KEY).(xid.ID); ok && !id.IsZero() {
		fields = append(fields, "reqTraceId", id)
	}

	if id, ok := ctx.Value(XLOG_REQ_TRACE_ID_CTX_KEY).(string); ok && id != "" {
		fields = append(fields, "reqTraceId", id)
	}

	fields = append(fields,
		"querySql", sqlStr,
		"queryRowsAffected", rows,
		"queryStartTime", begin,
		"queryEndTime", time.Now(),
		"queryDuration", FormatDuration(elapsed),
		"queryDurationNs", elapsed.Nanoseconds(),
	)

	// Mark slow queries
	if g.SlowQuery > 0 && elapsed > g.SlowQuery && g.LogLevel >= logger.Warn {
		fields = append(fields, "slow_query", true)
		g.Log.Debug(ctx, "slow query detected", fields...)
		return
	}

	if err != nil && g.LogLevel >= logger.Error {
		fields = append(fields, "err", err)
		g.Log.Debug(ctx, "executing query failed", fields...)
		return
	}

	if g.LogLevel >= logger.Info {
		g.Log.Debug(ctx, "executing query success", fields...)
	}
}
