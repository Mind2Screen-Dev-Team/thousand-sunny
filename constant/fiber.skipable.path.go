package constant

var (
	FiberSkipablePathFromMiddleware = [...]string{
		"/favicon.ico", "/openapi.json", "/openapi.yaml",
		"/docs", "/schemas", "/monitor",
	}
)
