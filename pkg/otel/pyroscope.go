package otel

import (
	"github.com/go-list-templ/users-service/pkg/config"
	"github.com/grafana/pyroscope-go"
)

type Pyroscope struct {
	*pyroscope.Profiler
}

func NewPyroscope(cfg *config.Config) (*Pyroscope, error) {
	profiler, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: cfg.App.Name,
		ServerAddress:   cfg.Otel.Endpoint,
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
	})
	if err != nil {
		return nil, err
	}

	return &Pyroscope{profiler}, err
}

func (p *Pyroscope) Shutdown() error {
	return p.Stop()
}
