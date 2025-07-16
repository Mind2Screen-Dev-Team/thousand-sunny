package dependency

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"resty.dev/v3"
)

func ProvideHttpClient(debuglog *xlog.DebugLogger) *resty.Client {
	var (
		xl = xlog.NewLogger(debuglog.Logger)
		l  = xlog.NewRestyV3Logger(xl)
		c  = resty.New().EnableTrace().EnableDebug()
	)

	return c.
		SetLogger(l).
		SetDebugLogFormatter(
			resty.DebugLogJSONFormatter,
		)
}
