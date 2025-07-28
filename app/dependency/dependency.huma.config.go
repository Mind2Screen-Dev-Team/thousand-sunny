package dependency

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/infra/http/middleware"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhuma"
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func ProvideHumaConfig(s config.Server) huma.Config {
	var (
		schemaPrefix = "#/components/schemas/"
		schemasPath  = "/schemas"
		registry     = huma.NewMapRegistry(schemaPrefix, huma.DefaultSchemaNamer)
		oapi         = huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info:    s.OAPI.Info,
			Components: &huma.Components{
				Schemas: registry,
			},
			Servers: s.OAPI.Server,
		}
	)

	// Register Fiber Monitor Metric into OAPI
	middleware.RegisterMiddlewareMonitorOAPI(&oapi)

	return huma.Config{
		OpenAPI:       &oapi,
		OpenAPIPath:   "/openapi",
		DocsPath:      "/docs",
		SchemasPath:   schemasPath,
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
		CreateHooks: []func(huma.Config) huma.Config{
			func(c huma.Config) huma.Config {
				c.Transformers = append(c.Transformers, xhuma.TransformerTraceIdSetter)
				return c
			},
		},
	}
}
