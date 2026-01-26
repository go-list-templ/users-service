package server

import (
	"context"
	"net"
	"net/http"

	"github.com/go-list-templ/grpc/pkg/config"
)

type HTTP struct {
	server http.Server
	config *config.Server
	errors chan error
}

func NewHTTP(cfg *config.Server) *HTTP {
	return &HTTP{
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

func (s *HTTP) Notify() <-chan error {
	return s.errors
}

func (s *HTTP) Start() {
	go func() {
		s.errors <- s.server.ListenAndServe()
		close(s.errors)
	}()
}

func (s *HTTP) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
