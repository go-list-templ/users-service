package otel

import (
	"context"
	"os"

	"github.com/go-list-templ/grpc/pkg/config"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.39.0"
)

type Telemetry struct {
	LoggerProvider *Logger
	TracerProvider *Trace
	MetricProvider *Metric
}

func NewTelemetry(cfg *config.Config) (*Telemetry, error) {
	ctx := context.Background()

	hostName, _ := os.Hostname()

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.App.Name),
		semconv.ServiceVersion(cfg.App.Version),
		semconv.HostName(hostName),
	)

	loggerProvider, err := NewLogger(ctx, res, &cfg.Otel)
	if err != nil {
		return nil, err
	}

	tracerProvider, err := NewTrace(ctx, res, &cfg.Otel)
	if err != nil {
		return nil, err
	}

	metricProvider, err := NewMetric(ctx, res, &cfg.Otel)
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
	var err error

	err = t.MetricProvider.Shutdown(ctx)
	err = t.TracerProvider.Shutdown(ctx)
	err = t.LoggerProvider.Shutdown(ctx)

	return err
}
