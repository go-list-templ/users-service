package grpc

import (
	pbgrpc "google.golang.org/grpc"

	"github.com/go-list-templ/grpc/internal/controller/grpc/v1/user"
	"github.com/go-list-templ/grpc/internal/usecase/user/command"
	"github.com/go-list-templ/grpc/internal/usecase/user/query"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"
)

type Router struct {
	server *pbgrpc.Server
	logger *zap.Logger
}

func NewRouter(s *pbgrpc.Server, l *zap.Logger) *Router {
	return &Router{server: s, logger: l}
}

func (r *Router) Release() {
	reflection.Register(r.server)
}

func (r *Router) User(uq *query.UserUsecase, uc *command.UserUsecase) {
	user.NewQueryService(r.server, uq, r.logger)
	user.NewCommandService(r.server, uc, r.logger)
}
