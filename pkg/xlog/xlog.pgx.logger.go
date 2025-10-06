package xlog

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/xid"
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
	Log Logger
}

// TraceQueryStart logs the start of a query.
func (t *PgxLogger) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	var (
		n      = time.Now()
		fields = make([]any, 0)
		id, _  = ctx.Value(XLOG_REQ_TRACE_ID_CTX_KEY).(xid.ID)
	)

	if !id.IsZero() {
		fields = append(fields, "reqTraceId", id)
	}

	fields = append(fields,
		"reqTraceId", id,
		"querySql", data.SQL,
		"queryArgs", data.Args,
		"queryStartTime", n,
	)
	t.Log.Debug(ctx, "start executing query", fields...)

	ctx = context.WithValue(ctx, queryStartTimeKey{}, n)
	ctx = context.WithValue(ctx, querySQLDataKey{}, data)
	return ctx
}

// TraceQueryEnd logs the end of a query.
func (t *PgxLogger) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	var (
		queryEndTime      = time.Now()
		fields            = make([]any, 0)
		queryStartTime, _ = ctx.Value(queryStartTimeKey{}).(time.Time) // Retrieve the start time from context
		queryDurr         = queryEndTime.Sub(queryStartTime)

		id, _        = ctx.Value(XLOG_REQ_TRACE_ID_CTX_KEY).(xid.ID)
		queryName, _ = ctx.Value(queryNameKey{}).(string)
		queryData, _ = ctx.Value(querySQLDataKey{}).(pgx.TraceQueryStartData) // Retrieve the trace query start data from context
		queryType    = getQueryType(data.CommandTag)                          // Determine query typex
	)

	if data.Err != nil {
		if !id.IsZero() {
			fields = append(fields, "reqTraceId", id)
		}

		if queryName != "" {
			fields = append(fields, "queryName", queryName)
		}

		if queryType != "" {
			fields = append(fields, "queryType", queryType)
		}

		fields = append(fields,
			"err", data.Err,
			"querySql", queryData.SQL,
			"queryArgs", queryData.Args,
			"queryEndTime", queryEndTime,
			"queryStartTime", queryStartTime,
			"queryDuration", FormatDuration(queryDurr),
			"queryDurationNs", queryDurr.Nanoseconds(),
		)
		t.Log.Debug(ctx, "executing query is failed", fields...)
		return
	}

	if !id.IsZero() {
		fields = append(fields, "reqTraceId", id)
	}

	if queryName != "" {
		fields = append(fields, "queryName", queryName)
	}

	if queryType != "" {
		fields = append(fields, "queryType", queryType)
	}

	fields = append(fields,
		"querySql", queryData.SQL,
		"queryArgs", queryData.Args,
		"queryRowsAffected", data.CommandTag.RowsAffected(),
		"queryEndTime", queryEndTime,
		"queryStartTime", queryStartTime,
		"queryDuration", FormatDuration(queryDurr),
		"queryDurationNs", queryDurr.Nanoseconds(),
	)
	t.Log.Debug(ctx, "executing query is success", fields...)
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
