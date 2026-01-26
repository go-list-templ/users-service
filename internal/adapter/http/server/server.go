package server

import (
	"context"
	"net"
	"net/http"

	"github.com/go-list-templ/grpc/config"
)

type Health struct {
	server http.Server
	config *config.Server
	errors chan error
}

func NewHealth(cfg *config.Server) *Health {
	return &Health{
		server: http.Server{
			Addr:              net.JoinHostPort("", cfg.HealthPort),
			Handler:           nil,
			ReadHeaderTimeout: cfg.HTTPTimeout,
			IdleTimeout:       cfg.IdleTimeout,
		},
		config: cfg,
		errors: make(chan error, 1),
	}
}

func (s *Health) Notify() <-chan error {
	return s.errors
}

func (s *Health) Start() {
	go func() {
		s.errors <- s.server.ListenAndServe()
		close(s.errors)
	}()
}

func (s *Health) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
