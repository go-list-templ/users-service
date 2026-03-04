package otel

import (
	"context"

	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"

	"github.com/go-list-templ/grpc/pkg/config"
	"go.opentelemetry.io/otel/sdk/resource"
)

type Telemetry struct {
	LoggerProvider *Logger
	TracerProvider *Trace
	MetricProvider *Metric
}

func NewTelemetry(cfg *config.Config) (*Telemetry, error) {
	ctx := context.Background()

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.App.Name),
		semconv.ServiceVersion(cfg.App.Version),
	)

	metricProvider, err := NewMetric(ctx, res, &cfg.Otel)
	if err != nil {
		return nil, err
	}

	tracerProvider, err := NewTrace(ctx, res, &cfg.Otel)
	if err != nil {
		return nil, err
	}

	loggerProvider, err := NewLogger(ctx, res, &cfg.Otel)
	if err != nil {
		return nil, err
	}

	return &Telemetry{
		LoggerProvider: loggerProvider,
		TracerProvider: tracerProvider,
		MetricProvider: metricProvider,
	}, nil
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	err := t.MetricProvider.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = t.TracerProvider.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = t.LoggerProvider.Shutdown(ctx)
	if err != nil {
		return err
	}

	return err
}
