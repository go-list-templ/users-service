package grpc

import (
	"github.com/go-list-templ/grpc/internal/controller/grpc/v1/user"
	"github.com/go-list-templ/grpc/internal/usecase/user/command"
	"github.com/go-list-templ/grpc/internal/usecase/user/query"
	"go.uber.org/zap"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewRouter(app *pbgrpc.Server, uq *query.UserUsecase, uc *command.UserUsecase, l *zap.Logger) {
	{
		user.NewQueryService(app, uq, l)
		user.NewCommandService(app, uc, l)
	}

	reflection.Register(app)
}
