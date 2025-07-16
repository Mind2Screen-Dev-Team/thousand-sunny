package dependency

import (
	"context"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"resty.dev/v3"
)

func ProvideHttpClient(debug *xlog.DebugLogger) *resty.Client {
	var (
		logger = xlog.NewLogger(debug.Logger)
		client = resty.New().EnableTrace().EnableDebug()
	)

	return client.OnDebugLog(func(dl *resty.DebugLog) {
		var (
			fields = []any{
				"data", dl,
			}
			actionMsg = dl.Request.Header.Get("X-Action-Msg")
			tid       = dl.Request.Header.Get("X-Trace-ID")
		)
		if actionMsg == "" {
			actionMsg = "http resty debug"
		}
		if tid != "" {
			fields = append(fields, "traceId", tid)
		}

		logger.Info(context.Background(), actionMsg, fields...)
	})
}
