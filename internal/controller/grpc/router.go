package grpc

import (
	v1 "github.com/go-list-templ/grpc/internal/controller/grpc/v1"
	"github.com/go-list-templ/grpc/internal/usecase"
	"go.uber.org/zap"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewRouter(app *pbgrpc.Server, u *usecase.User, l *zap.Logger) {
	{
		v1.NewUserRoute(app, u, l)
	}

	reflection.Register(app)
}
