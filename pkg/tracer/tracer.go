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
	Endpoint string
}

type AppTracer struct {
	TraceProvider *sdkTrace.TracerProvider
	Tracer        trace.Tracer
}

func NewAppTracer(ctx context.Context, cfg *Config, serviceName, version string) (*AppTracer, error) {
	// Validate configuration
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is missing in the tracer configuration")
	}

	// Initialize Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Endpoint)))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Jaeger exporter: %w", err)
	}

	// Create resource
	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
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

	// Create span processor
	bsp := sdkTrace.NewBatchSpanProcessor(exp)

	// Create TracerProvider
	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithSpanProcessor(bsp),
		sdkTrace.WithResource(res),
	)

	// Create AppTracer instance
	at := &AppTracer{
		TraceProvider: tp,
	}
	at.Tracer = tp.Tracer("")

	// Set TextMap propagator and TracerProvider globally
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tp)

	// Log successful tracer initialization
	fmt.Println("Tracer initialized successfully. Endpoint:", cfg.Endpoint)

	return at, nil
}
