package server

import (
	"context"
	"net"

	"github.com/go-list-templ/grpc/pkg/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type GRPC struct {
	Server *grpc.Server
	ctx    context.Context
	eg     *errgroup.Group
	config *config.Server
	errors chan error
}

func New(cfg *config.Server) *GRPC {
	ka := keepalive.ServerParameters{
		MaxConnectionIdle: cfg.GRPCMaxConnIdle,
		MaxConnectionAge:  cfg.GRPCMaxConnAge,
		Time:              cfg.GRPCTime,
		Timeout:           cfg.GRPCTimeout,
	}

	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.KeepaliveParams(ka),
		grpc.ConnectionTimeout(cfg.GRPCTimeout),
	)

	return &GRPC{
		Server: server,
		ctx:    context.Background(),
		eg:     &errgroup.Group{},
		config: cfg,
		errors: make(chan error, 1),
	}
}

func (s *GRPC) Notify() <-chan error {
	return s.errors
}

func (s *GRPC) Start() {
	s.eg.Go(func() error {
		var lc net.ListenConfig

		ln, err := lc.Listen(s.ctx, "tcp", net.JoinHostPort("", s.config.GRPCPort))
		if err != nil {
			s.errors <- err

			close(s.errors)

			return err
		}

		err = s.Server.Serve(ln)
		if err != nil {
			s.errors <- err

			close(s.errors)

			return err
		}

		return nil
	})
}

func (s *GRPC) Stop() {
	s.Server.GracefulStop()
}
