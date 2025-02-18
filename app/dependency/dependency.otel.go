package dependency

import (
	"context"
	"fmt"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xtracer"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func ProvideOtelConfig(c config.Cfg, s config.Server) xtracer.Config {
	return xtracer.Config{
		Tracer: c.Otel.Tracer,
		Metric: c.Otel.Metric,

		ModuleName:    fmt.Sprintf("%s/%s", c.App.Project, s.Name),
		ServerName:    fmt.Sprintf("%s/%s", c.App.Project, s.Name),
		ServerAddress: s.Address,

		GrpcHost: c.Otel.Server.GrpcHost,
		GrpcPort: c.Otel.Server.GrpcPort,
	}
}

func ProvideOtelResource(c xtracer.Config) (*resource.Resource, error) {
	res, err := xtracer.NewResource(context.Background(), c)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ProvideOtelGrpcClient(c xtracer.Config) (*xtracer.GrpcClient, error) {
	client, err := xtracer.NewGrpcClient(c)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type OtelParamFx struct {
	fx.In

	Lc         fx.Lifecycle
	Cfg        xtracer.Config
	GrpcClient *xtracer.GrpcClient
	Resource   *resource.Resource
}

func ProvideOtelTracer(ctx context.Context, p OtelParamFx) (trace.Tracer, error) {
	tracer, shutdownFn, err := xtracer.NewOtelTracer(ctx, p.Cfg, p.Resource, p.GrpcClient.ClientConn)
	if err != nil {
		return nil, err
	}

	p.Lc.Append(fx.Hook{OnStop: shutdownFn})

	return tracer, nil
}

func ProvideOtelMetric(ctx context.Context, p OtelParamFx) (metric.Meter, error) {
	meter, shutdownFn, err := xtracer.NewOtelMeter(
		ctx,
		p.Cfg,
		p.Resource,
		p.GrpcClient.ClientConn,
	)
	if err != nil {
		return nil, err
	}

	p.Lc.Append(fx.Hook{OnStop: shutdownFn})

	return meter, nil
}
