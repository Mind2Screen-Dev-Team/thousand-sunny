package xlog

import (
	"context"
	"fmt"
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

func PgxQueryName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, queryNameKey{}, name)
}

// PgxLogger is a custom QueryTracer implementation for pgx using zerolog.
type PgxLogger struct {
	Log zerolog.Logger
}

// TraceQueryStart logs the start of a query.
func (t *PgxLogger) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	var (
		n     = time.Now()
		e     = t.Log.Info()
		id, _ = ctx.Value(XLOG_TRACE_ID_CTX_KEY).(xid.ID)
	)

	if !id.IsZero() {
		e = e.Str("req.trace.id", id.String())
	}

	e = e.Str("query.sql", data.SQL)
	e = e.Any("query.args", data.Args)
	e = e.Time("query.start.time", n)
	e.Msg("start executing query")

	ctx = context.WithValue(ctx, queryStartTimeKey{}, n)
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
		e = e.Str("query.duration", FormatDuration(queryDurr))
		e = e.Int64("query.duration.ns", queryDurr.Nanoseconds())
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
	e = e.Str("query.duration", FormatDuration(queryDurr))
	e = e.Int64("query.duration.ns", queryDurr.Nanoseconds())
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

// FormatDuration dynamically formats time.Duration into ns, µs, ms, s, m, h, or d
func FormatDuration(d time.Duration) string {
	const day = 24 * time.Hour

	switch {
	case d >= day:
		return fmt.Sprintf("%.2f d", d.Hours()/24)
	case d >= time.Hour:
		return fmt.Sprintf("%.2f h", d.Hours())
	case d >= time.Minute:
		return fmt.Sprintf("%.2f m", d.Minutes())
	case d >= time.Second:
		return fmt.Sprintf("%.3f s", d.Seconds())
	case d >= time.Millisecond:
		return fmt.Sprintf("%.3f ms", float64(d)/float64(time.Millisecond))
	case d >= time.Microsecond:
		return fmt.Sprintf("%.3f µs", float64(d)/float64(time.Microsecond))
	default:
		return fmt.Sprintf("%d ns", d.Nanoseconds())
	}
}
