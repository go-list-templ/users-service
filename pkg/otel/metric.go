package otel

import (
	"context"
	"fmt"
	"time"

	otelmetric "go.opentelemetry.io/otel/metric"

	"github.com/go-list-templ/users-service/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

const MetricInterval = time.Second * 10

type Metric struct {
	otelmetric.Meter

	Provider *metric.MeterProvider
}

func NewMetricProvider(ctx context.Context, res *resource.Resource, cfg *config.Otel) (*metric.MeterProvider, error) {
	options := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.Endpoint),
	}

	if !cfg.IsTLS {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(MetricInterval))),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(provider)

	return provider, nil
}

func NewMetric(ctx context.Context, res *resource.Resource, cfg *config.Config) (*Metric, error) {
	provider, err := NewMetricProvider(ctx, res, &cfg.Otel)
	if err != nil {
		return nil, fmt.Errorf("failed to create meter: %w", err)
	}

	return &Metric{provider.Meter(cfg.App.Name), provider}, nil
}

func (m *Metric) Shutdown(ctx context.Context) error {
	return m.Provider.Shutdown(ctx)
}
