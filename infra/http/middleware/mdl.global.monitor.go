package middleware

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func ProvideMonitor() Monitor {
	return Monitor{}
}

type Monitor struct{}

func (Monitor) Name() string {
	return "monitor"
}

func (Monitor) App(app *fiber.App) {
	app.Get("/monitor", monitor.New(monitor.ConfigDefault))
}

func (Monitor) Serve(c *fiber.Ctx) error {
	return c.Next()
}

// FiberMonitorMetricStats represents monitoring metrics for both process and operating system level.
type FiberMonitorMetricStats struct {
	PID FiberMonitorMetricStatsPID `json:"pid" doc:"Process-level metrics" example:"{\"cpu\":0.050654302209225406,\"ram\":24383488,\"conns\":17}"`
	OS  FiberMonitorMetricStatsOS  `json:"os" doc:"Operating system-level metrics" example:"{\"cpu\":23.7611181702772,\"ram\":11562934272,\"total_ram\":17179869184,\"load_avg\":2.91064453125,\"conns\":338}"`
}

// FiberMonitorMetricStatsPID represents the metrics for the current process.
type FiberMonitorMetricStatsPID struct {
	CPU   float64 `json:"cpu" doc:"CPU usage percentage of the process" example:"0.050654302209225406"`
	RAM   uint64  `json:"ram" doc:"RAM usage by the process in bytes" example:"24383488"`
	Conns int     `json:"conns" doc:"Active connections handled by the process" example:"17"`
}

// FiberMonitorMetricStatsOS represents the metrics for the entire operating system.
type FiberMonitorMetricStatsOS struct {
	CPU      float64 `json:"cpu" doc:"Total CPU usage percentage of the OS" example:"23.7611181702772"`
	RAM      uint64  `json:"ram" doc:"Used RAM in bytes" example:"11562934272"`
	TotalRAM uint64  `json:"total_ram" doc:"Total available RAM in bytes" example:"17179869184"`
	LoadAvg  float64 `json:"load_avg" doc:"System load average" example:"2.91064453125"`
	Conns    int     `json:"conns" doc:"Total active network connections in the OS" example:"338"`
}

func RegisterMiddlewareMonitorOAPI(oapi *huma.OpenAPI) {
	oapi.AddOperation(&huma.Operation{
		OperationID:   "fiber-monitor",
		Method:        http.MethodGet,
		Path:          "/monitor",
		Summary:       "Fiber Monitoring API",
		Description:   "Retrieves monitoring data for the application and underlying system. This includes metrics such as CPU usage, memory consumption, number of active connections, and system load average. Useful for health checks, diagnostics, and observability dashboards. For fiber monitoring page, visit the url 'api_base_url/monitor'",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Metrics"},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Description: "Successful metrics response",
				Content: map[string]*huma.MediaType{
					fiber.MIMEApplicationJSON: {
						Schema: huma.SchemaFromType(oapi.Components.Schemas, reflect.TypeOf(FiberMonitorMetricStats{})),
						Example: FiberMonitorMetricStats{
							PID: FiberMonitorMetricStatsPID{
								CPU:   0.025684832030902268,
								RAM:   24002560,
								Conns: 26,
							},
							OS: FiberMonitorMetricStatsOS{
								CPU:      14.145838701382443,
								RAM:      11910316032,
								TotalRAM: 17179869184,
								LoadAvg:  7.14501953125,
								Conns:    401,
							},
						},
					},
				},
			},
		},
	})
}
