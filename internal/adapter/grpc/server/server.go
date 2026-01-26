package server

import (
	"net"

	"github.com/go-list-templ/grpc/pkg/config"
	"google.golang.org/grpc"
)

type GRPC struct {
	Server *grpc.Server
	config *config.Server
	errors chan error
}

func New(cfg *config.Server) *GRPC {
	grpcServer := grpc.NewServer()

	return &GRPC{
		Server: grpcServer,
		config: cfg,
		errors: make(chan error, 1),
	}
}

func (s *GRPC) Notify() <-chan error {
	return s.errors
}

func (s *GRPC) Start() {
	go func() {
		lis, err := net.Listen("tcp", net.JoinHostPort("", s.config.GRPCPort))
		if err != nil {
			s.errors <- err
		}

		s.errors <- s.Server.Serve(lis)
		close(s.errors)
	}()
}

func (s *GRPC) Stop() {
	s.Server.GracefulStop()
}
