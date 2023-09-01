package probes

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Config struct {
	Path      string `mapstructure:"PATH"`
	Port      string `mapstructure:"PORT"`
	Prefix    string `mapstructure:"PREFIX"`
	MeterName string `mapstructure:"METER_NAME"`
}

type AppMetric struct {
	MetricProvider *sdkMetric.MeterProvider
	Meter          metric.Meter
}

func NewAppMetric(ctx context.Context, cfg *Config, serviceName, version string) (*AppMetric, error) {
	exporter, err := prometheus.New(prometheus.WithNamespace(cfg.Prefix))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Prometheus exporter: %w", err)
	}

	resource, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(version),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	metricProvider := sdkMetric.NewMeterProvider(
		sdkMetric.WithResource(resource),
		sdkMetric.WithReader(exporter))

	ap := &AppMetric{
		MetricProvider: metricProvider,
		Meter:          metricProvider.Meter(cfg.MeterName),
	}

	otel.SetMeterProvider(metricProvider)

	return ap, err
}
