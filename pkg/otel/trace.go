package otel

import (
	"context"
	"fmt"

	"github.com/go-list-templ/users-service/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Trace struct {
	Provider *trace.TracerProvider
}

func NewTrace(ctx context.Context, res *resource.Resource, cfg *config.Otel) (*Trace, error) {
	provider, err := NewTraceProvider(ctx, res, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer: %w", err)
	}

	return &Trace{provider}, nil
}

func NewTraceProvider(ctx context.Context, res *resource.Resource, cfg *config.Otel) (*trace.TracerProvider, error) {
	options := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	}

	if !cfg.IsTLS {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(provider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return provider, nil
}

func (t *Trace) Shutdown(ctx context.Context) error {
	return t.Provider.Shutdown(ctx)
}
