package otel

import (
	"context"
	"fmt"
	"os"

	"github.com/go-list-templ/grpc/pkg/config"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger

	provider *log.LoggerProvider
}

func NewLogger(ctx context.Context, res *resource.Resource, cfg *config.Otel) (*Logger, error) {
	provider, err := NewLoggerProvider(ctx, res, cfg)
	if err != nil {
		return nil, err
	}

	logger := zap.New(
		zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(os.Stdout),
				zapcore.InfoLevel,
			),
			otelzap.NewCore("name", otelzap.WithLoggerProvider(provider)),
		),
	)

	return &Logger{logger, provider}, nil
}

// NewLoggerProvider todo delete or auto set WithInsecure()
func NewLoggerProvider(ctx context.Context, res *resource.Resource, cfg *config.Otel) (*log.LoggerProvider, error) {
	exporter, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint(cfg.Endpoint),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}

	processor := log.NewBatchProcessor(exporter)

	provider := log.NewLoggerProvider(
		log.WithProcessor(processor),
		log.WithResource(res),
	)

	return provider, nil
}

func (t *Logger) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}
