package otel

import (
	"os"

	"github.com/go-list-templ/grpc/pkg/config"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.39.0"
)

func NewResource(cfg *config.App) *resource.Resource {
	hostName, _ := os.Hostname()

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.Name),
		semconv.ServiceVersion(cfg.Version),
		semconv.HostName(hostName),
	)
}
