package xlog

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

type (
	queryNameKey      struct{}
	queryStartTimeKey struct{}
	querySQLDataKey   struct{}
)

// Define context keys to store start time
func QueryName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, queryNameKey{}, name)
}

// PgxLogger is a custom QueryTracer implementation for pgx using zerolog.
type PgxLogger struct {
	Log zerolog.Logger
}

// TraceQueryStart logs the start of a query.
func (t *PgxLogger) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	var (
		queryStartTime = time.Now()
		id, _          = ctx.Value(XLOG_TRACE_ID_CTX_KEY).(xid.ID)
		e              = t.Log.Info()
	)

	if !id.IsZero() {
		e = e.Str("req.trace.id", id.String())
	}

	e = e.Str("query.sql", data.SQL)
	e = e.Any("query.args", data.Args)
	e = e.Time("query.start.time", queryStartTime)
	e.Msg("start executing query")

	ctx = context.WithValue(ctx, queryStartTimeKey{}, queryStartTime)
	ctx = context.WithValue(ctx, querySQLDataKey{}, data)
	return ctx
}

// TraceQueryEnd logs the end of a query.
func (t *PgxLogger) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	var (
		queryEndTime      = time.Now()
		queryStartTime, _ = ctx.Value(queryStartTimeKey{}).(time.Time) // Retrieve the start time from context
		queryDurr         = queryEndTime.Sub(queryStartTime)

		id, _        = ctx.Value(XLOG_TRACE_ID_CTX_KEY).(xid.ID)
		queryName, _ = ctx.Value(queryNameKey{}).(string)
		queryData, _ = ctx.Value(querySQLDataKey{}).(pgx.TraceQueryStartData) // Retrieve the trace query start data from context
		queryType    = getQueryType(data.CommandTag)                          // Determine query typex
	)

	if data.Err != nil {
		e := t.Log.Error()

		if !id.IsZero() {
			e = e.Str("req.trace.id", id.String())
		}

		if queryName != "" {
			e = e.Str("query.name", queryName)
		}

		if queryType != "" {
			e = e.Str("query.type", queryType)
		}

		e = e.Err(data.Err)
		e = e.Str("query.sql", queryData.SQL)
		e = e.Any("query.args", queryData.Args)
		e = e.Time("query.end.time", queryEndTime)
		e = e.Time("query.start.time", queryStartTime)
		e = e.Dur("query.duration", queryDurr)
		e.Msg("executing query is failed")
		return
	}

	// Log execution success with rows affected
	e := t.Log.Info()
	if !id.IsZero() {
		e = e.Str("req.trace.id", id.String())
	}
	if queryName != "" {
		e = e.Str("query.name", queryName)
	}
	if queryType != "" {
		e = e.Str("query.type", queryType)
	}
	e = e.Str("query.sql", queryData.SQL)
	e = e.Any("query.args", queryData.Args)
	e = e.Int64("query.rows.affected", data.CommandTag.RowsAffected())
	e = e.Time("query.end.time", queryEndTime)
	e = e.Time("query.start.time", queryStartTime)
	e = e.Dur("query.duration", queryDurr)
	e.Msg("executing query is success")
}

// getQueryType determines the type of SQL command from CommandTag.
func getQueryType(ct pgconn.CommandTag) string {
	switch {
	case ct.Select():
		return "SELECT"
	case ct.Insert():
		return "INSERT"
	case ct.Update():
		return "UPDATE"
	case ct.Delete():
		return "DELETE"
	default:
		return ""
	}
}
