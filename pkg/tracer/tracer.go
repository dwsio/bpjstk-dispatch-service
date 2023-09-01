package tracer

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Endpoint string `mapstructure:"ENDPOINT"`
	Prefix   string `mapstructure:"PREFIX"`
}

type AppTracer struct {
	TraceProvider *sdkTrace.TracerProvider
	Tracer        trace.Tracer
}

func NewAppTracer(ctx context.Context, cfg *Config, serviceName, version string) (*AppTracer, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is missing in the tracer configuration")
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Endpoint)))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Jaeger exporter: %w", err)
	}

	resource, err := resource.New(
		ctx,
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(version),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	batchSpanProcessor := sdkTrace.NewBatchSpanProcessor(exporter)

	tracerProvider := sdkTrace.NewTracerProvider(
		sdkTrace.WithSpanProcessor(batchSpanProcessor),
		sdkTrace.WithResource(resource),
	)

	at := &AppTracer{
		TraceProvider: tracerProvider,
		Tracer:        tracerProvider.Tracer(cfg.Prefix),
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	fmt.Println("Tracer initialized successfully. Endpoint:", cfg.Endpoint)

	return at, nil
}
