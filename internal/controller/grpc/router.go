package grpc

import (
	v1 "github.com/go-list-templ/grpc/internal/controller/grpc/v1"
	pbgrpc "google.golang.org/grpc"

	"github.com/go-list-templ/grpc/internal/usecase/user"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"
)

func NewRouter(s *pbgrpc.Server, u *user.Usecase, l *zap.Logger) {
	{
		v1.NewUserRoute(s, u, l)
	}

	reflection.Register(s)
}
