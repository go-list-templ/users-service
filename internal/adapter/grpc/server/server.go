package server

import (
	"context"
	"net"

	"github.com/go-list-templ/grpc/pkg/config"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type GRPC struct {
	Server *grpc.Server
	ctx    context.Context
	eg     *errgroup.Group
	config *config.Server
	errors chan error
}

func New(cfg *config.Server) *GRPC {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1)

	return &GRPC{
		Server: grpc.NewServer(),
		ctx:    ctx,
		eg:     group,
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

		ln, err := lc.Listen(s.ctx, "tcp", s.config.GRPCPort)
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
