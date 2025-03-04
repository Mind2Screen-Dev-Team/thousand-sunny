package xtracer

import (
	"context"
	"fmt"
	"time"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/rs/xid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func Start(tracer trace.Tracer, ctx context.Context, span string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if reqTraceId, ok := ctx.Value(xlog.XLOG_REQ_TRACE_ID_CTX_KEY).(xid.ID); ok {
		opts = append(opts, trace.WithAttributes(attribute.String("req_trace_id", reqTraceId.String())))
	}

	opts = append(opts, trace.WithTimestamp(time.Now()))

	return tracer.Start(ctx, span, opts...)
}

// # Configuration

type Config struct {
	Tracer bool
	Metric bool
	Logs   bool

	ModuleName    string
	ServerName    string
	ServerAddress string

	GrpcHost string
	GrpcPort int
}

type (
	ShutdownFn = func(context.Context) error
	GrpcClient struct{ *grpc.ClientConn }
)

func NewGrpcClient(cfg Config) (*GrpcClient, error) {
	if !cfg.Tracer && !cfg.Metric {
		return nil, nil
	}

	// It connects the OpenTelemetry Collector through local gRPC connection.
	// You may replace `localhost:4317` with your endpoint.
	conn, err := grpc.NewClient(
		// target url
		fmt.Sprintf("%s:%d", cfg.GrpcHost, cfg.GrpcPort),
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("'failed to create otel gRPC connection to collector': %w", err)
	}

	return &GrpcClient{conn}, err
}

func NewResource(ctx context.Context, cfg Config) (*sdkresource.Resource, error) {
	if !cfg.Tracer && !cfg.Metric {
		return nil, nil
	}

	res, err := sdkresource.New(ctx,
		sdkresource.WithAttributes(
			// The service name used to display traces in backends
			semconv.ServiceNameKey.String(cfg.ServerName),
		),
	)

	return res, err
}

func NewOtelTracer(ctx context.Context, cfg Config, res *sdkresource.Resource, conn *grpc.ClientConn) (trace.Tracer, ShutdownFn, error) {
	var (
		tracerShutdownFn ShutdownFn
		tracerProvider   trace.TracerProvider
	)

	if cfg.Tracer {
		// Set up a trace exporter
		traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, nil, fmt.Errorf("'failed to create otel trace exporter': %w", err)
		}

		// Register the trace exporter with a TracerProvider, using a batch
		// span processor to aggregate spans before export.
		sp := sdktrace.NewBatchSpanProcessor(traceExporter)
		tracerProvider = sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(res),
			sdktrace.WithSpanProcessor(sp),
		)
		tracerShutdownFn = traceExporter.Shutdown
	} else {
		tracerProvider = tracenoop.NewTracerProvider()
		tracerShutdownFn = func(ctx context.Context) error { return nil }
	}
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tracer := otel.Tracer(cfg.ModuleName,
		trace.WithInstrumentationAttributes(
			attribute.String("server_name", cfg.ServerName),
			attribute.String("server_addr", cfg.ServerAddress),
		),
	)

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracer, tracerShutdownFn, nil
}

func NewOtelMeter(ctx context.Context, cfg Config, res *sdkresource.Resource, conn *grpc.ClientConn) (metric.Meter, ShutdownFn, error) {
	var (
		meterShutdownFn ShutdownFn
		meterProvider   metric.MeterProvider
	)

	if cfg.Metric {
		metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, nil, fmt.Errorf("'failed to create otel metrics exporter': %w", err)
		}

		meterProvider = sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
			sdkmetric.WithResource(res),
		)
		meterShutdownFn = metricExporter.Shutdown
	} else {
		meterProvider = metricnoop.NewMeterProvider()
		meterShutdownFn = func(ctx context.Context) error { return nil }
	}

	otel.SetMeterProvider(meterProvider)

	meter := otel.Meter(cfg.ModuleName,
		metric.WithInstrumentationAttributes(
			attribute.String("server_name", cfg.ServerName),
			attribute.String("server_addr", cfg.ServerAddress),
		),
	)

	return meter, meterShutdownFn, nil
}
