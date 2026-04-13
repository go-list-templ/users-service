package otel

import (
	"github.com/go-list-templ/users-service/pkg/config"
	otelpyroscope "github.com/grafana/otel-profiling-go"
	"github.com/grafana/pyroscope-go"
	"go.opentelemetry.io/otel"
	"time"
)

type Pyroscope struct {
	*pyroscope.Profiler
}

func NewPyroscope(cfg *config.Config, trace *Trace) (*Pyroscope, error) {
	otel.SetTracerProvider(otelpyroscope.NewTracerProvider(trace.Provider))

	profiler, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: cfg.App.Name,
		ServerAddress:   cfg.Otel.PyroscopeEndpoint,
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
		Tags: map[string]string{
			"version": cfg.App.Version,
		},
		UploadRate: 15 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &Pyroscope{profiler}, err
}

func (p *Pyroscope) Shutdown() error {
	return p.Stop()
}
