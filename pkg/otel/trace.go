package otel

import (
	"context"
	"fmt"

	"github.com/go-list-templ/grpc/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Trace struct {
	provider *trace.TracerProvider
}

func NewTrace(ctx context.Context, res *resource.Resource, cfg *config.Otel) (*Trace, error) {
	provider, err := NewTraceProvider(ctx, res, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer: %w", err)
	}

	return &Trace{provider}, nil
}

// NewTraceProvider todo delete or auto set WithInsecure()
func NewTraceProvider(ctx context.Context, res *resource.Resource, cfg *config.Otel) (*trace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(provider)

	return provider, nil
}

func (t *Trace) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}
