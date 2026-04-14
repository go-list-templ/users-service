package otel

import (
	otelpyroscope "github.com/grafana/otel-profiling-go"

	"github.com/go-list-templ/users-service/pkg/config"
	"github.com/grafana/pyroscope-go"
	"go.opentelemetry.io/otel"
)

type Pyroscope struct {
	*pyroscope.Profiler
}

func NewPyroscope(cfg *config.Config, trace *Trace) (*Pyroscope, error) {
	profiler, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: cfg.App.Name,
		ServerAddress:   cfg.Otel.PyroscopeEndpoint,
	})
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(otelpyroscope.NewTracerProvider(trace.Provider))

	return &Pyroscope{profiler}, err
}

func (p *Pyroscope) Shutdown() error {
	return p.Stop()
}
