package otel

import (
	"context"
	"fmt"

	"github.com/go-list-templ/grpc/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

type Metric struct {
	provider *metric.MeterProvider
}

func NewMetric(ctx context.Context, res *resource.Resource, cfg *config.Otel) (*Metric, error) {
	provider, err := NewMetricProvider(ctx, res, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create meter: %w", err)
	}

	return &Metric{provider}, nil
}

// NewMetricProvider todo delete or auto set WithInsecure()
func NewMetricProvider(ctx context.Context, res *resource.Resource, cfg *config.Otel) (*metric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(cfg.Endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(provider)

	return provider, nil
}

func (t *Metric) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}
