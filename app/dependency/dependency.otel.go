package dependency

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xtracer"
)

func ProvideOtelConfig(c config.Cfg, s config.Server) xtracer.Config {
	return xtracer.Config{
		Tracer: c.Otel.Tracer,
		Metric: c.Otel.Metric,
		Logs:   c.Otel.Logs,

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

func ProvideOtelTracer(p OtelParamFx) (trace.Tracer, error) {
	tracer, shutdownFn, err := xtracer.NewOtelTracer(context.Background(), p.Cfg, p.Resource, p.GrpcClient.ClientConn)
	if err != nil {
		return nil, err
	}

	p.Lc.Append(fx.Hook{OnStop: shutdownFn})

	return tracer, nil
}

func ProvideOtelMetric(p OtelParamFx) (metric.Meter, error) {
	meter, shutdownFn, err := xtracer.NewOtelMeter(
		context.Background(),
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

func ProvideOtelLog(p OtelParamFx) (*log.LoggerProvider, error) {
	if !p.Cfg.Logs {
		return nil, nil
	}

	logExporter, err := otlploggrpc.New(
		context.Background(),
		otlploggrpc.WithGRPCConn(p.GrpcClient.ClientConn),
	)
	if err != nil {
		return nil, err
	}

	var (
		logBatchProcessor = log.NewBatchProcessor(
			logExporter,
			log.WithExportMaxBatchSize(512),
			log.WithExportInterval(2*time.Second),
		)

		logProvider = log.NewLoggerProvider(
			log.WithResource(p.Resource),
			log.WithProcessor(logBatchProcessor),
		)
	)

	global.SetLoggerProvider(logProvider)

	p.Lc.Append(fx.Hook{OnStop: func(ctx context.Context) error {
		var errs []error

		if err := logProvider.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
		if err := logBatchProcessor.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
		if err := logExporter.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}

		return nil
	}})

	return logProvider, nil
}
