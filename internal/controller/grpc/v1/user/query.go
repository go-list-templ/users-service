package user

import (
	"context"
	"fmt"

	v1 "github.com/go-list-templ/proto/gen/api/user/v1"
	pbgrpc "google.golang.org/grpc"

	"github.com/go-list-templ/grpc/internal/usecase/user/query"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type QueryService struct {
	v1.UserQueryServiceServer

	usecase *query.UserUsecase
	logger  *zap.Logger
}

func NewQueryService(server *pbgrpc.Server, u *query.UserUsecase, l *zap.Logger) {
	service := &QueryService{usecase: u, logger: l}
	{
		v1.RegisterUserQueryServiceServer(server, service)
	}
}

func (q *QueryService) AllUsers(ctx context.Context, _ *v1.AllUsersRequest) (*v1.AllUsersResponse, error) {
	allUsers, err := q.usecase.All(ctx)
	if err != nil {
		q.logger.Warn("grpc - v1 - AllUsers", zap.Any("error:", err.Error()))

		return nil, fmt.Errorf("grpc - v1 - AllUsers: %w", err)
	}

	users := make([]*v1.User, len(allUsers))

	for i, user := range allUsers {
		users[i] = &v1.User{
			Id:        user.ID.Value().String(),
			Name:      user.Name.Value(),
			Email:     user.Email.Value(),
			Avatar:    user.Avatar.Value(),
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		}
	}

	return &v1.AllUsersResponse{
		Users: users,
	}, nil
}
