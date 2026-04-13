package otel

import (
	"context"

	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"

	"github.com/go-list-templ/users-service/pkg/config"
	"go.opentelemetry.io/otel/sdk/resource"
)

type Telemetry struct {
	Pyroscope *Pyroscope
	Logger    *Logger
	Tracer    *Trace
	Metric    *Metric
}

func NewTelemetry(cfg *config.Config) (*Telemetry, error) {
	ctx := context.Background()

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.App.Name),
		semconv.ServiceVersion(cfg.App.Version),
	)

	metric, err := NewMetric(ctx, res, cfg)
	if err != nil {
		return nil, err
	}

	tracer, err := NewTrace(ctx, res, &cfg.Otel)
	if err != nil {
		return nil, err
	}

	logger, err := NewLogger(ctx, res, cfg)
	if err != nil {
		return nil, err
	}

	pyroscope, err := NewPyroscope(cfg, tracer)
	if err != nil {
		return nil, err
	}

	return &Telemetry{
		Pyroscope: pyroscope,
		Logger:    logger,
		Tracer:    tracer,
		Metric:    metric,
	}, nil
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	err := t.Metric.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = t.Tracer.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = t.Logger.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = t.Pyroscope.Shutdown()
	if err != nil {
		return err
	}

	return err
}
