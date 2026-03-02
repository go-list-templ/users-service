package otel

import (
	"context"
	"fmt"
	"os"

	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/go-list-templ/grpc/pkg/config"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type TelemetryProvider interface {
	TraceStart(context.Context, string) (context.Context, oteltrace.Span)
	Shutdown(context.Context)
}

// Telemetry is a wrapper around the OpenTelemetry logger, meter, and tracer.
type Telemetry struct {
	Logger *zap.Logger

	lp     *log.LoggerProvider
	mp     *metric.MeterProvider
	tp     *trace.TracerProvider
	tracer oteltrace.Tracer
}

// NewTelemetry creates a new telemetry instance.
func NewTelemetry(ctx context.Context, cfg *config.Config) (*Telemetry, error) {
	rp := newResource(&cfg.App)

	lp, err := newLoggerProvider(ctx, rp, &cfg.Otel)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	logger := zap.New(
		zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
			otelzap.NewCore("name", otelzap.WithLoggerProvider(lp)),
		),
	)

	mp, err := newMeterProvider(ctx, rp, &cfg.Otel)
	if err != nil {
		return nil, fmt.Errorf("failed to create meter: %w", err)
	}

	tp, err := newTracerProvider(ctx, rp, &cfg.Otel)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer: %w", err)
	}
	tracer := tp.Tracer("service-name-grpc")

	return &Telemetry{
		lp:     lp,
		mp:     mp,
		tp:     tp,
		Logger: logger,
		tracer: tracer,
	}, nil
}

// TraceStart starts a new span with the given name. The span must be ended by calling End.
func (t *Telemetry) TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span) {
	return t.tracer.Start(ctx, name)
}

// Shutdown shuts down the logger, meter, and tracer.
func (t *Telemetry) Shutdown(ctx context.Context) {
	t.lp.Shutdown(ctx)
	t.mp.Shutdown(ctx)
	t.tp.Shutdown(ctx)
}
