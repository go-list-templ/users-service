package server

import (
	"net"

	"github.com/go-list-templ/grpc/config"
	"google.golang.org/grpc"
)

type API struct {
	Server *grpc.Server
	config *config.Server
	errors chan error
}

func NewAPI(cfg *config.Server) *API {
	grpcServer := grpc.NewServer()

	return &API{
		Server: grpcServer,
		config: cfg,
		errors: make(chan error, 1),
	}
}

func (s *API) Notify() <-chan error {
	return s.errors
}

func (s *API) Start() {
	go func() {
		lis, err := net.Listen("tcp", net.JoinHostPort("", s.config.GRPCPort))
		if err != nil {
			s.errors <- err
		}

		s.errors <- s.Server.Serve(lis)
		close(s.errors)
	}()
}

func (s *API) Stop() {
	s.Server.GracefulStop()
}
