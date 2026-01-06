package v1

import (
	"context"
	"fmt"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/usecase"
	v1 "github.com/go-list-templ/proto/gen/api/user/v1"
	"go.uber.org/zap"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	v1.UserServiceServer

	usecase usecase.User
	logger  zap.Logger
}

func NewUserRoute(server *pbgrpc.Server, u usecase.User, l zap.Logger) {
	r := &UserService{usecase: u, logger: l}
	{
		v1.RegisterUserServiceServer(server, r)
	}
}

func (r *UserService) CreateUser(ctx context.Context, request *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	user, err := entity.NewUser(request.GetName(), request.GetEmail())
	if err != nil {
		r.logger.Warn("grpc - v1 - NewUser", zap.Any("error:", err.Error()))

		return nil, fmt.Errorf("grpc - v1 - NewUser: %w", err)
	}

	createdUser, err := r.usecase.Create(ctx, *user)
	if err != nil {
		r.logger.Warn("grpc - v1 - CreateUser", zap.Any("error:", err.Error()))

		return nil, fmt.Errorf("grpc - v1 - CreateUser: %w", err)
	}

	return &v1.CreateUserResponse{
		User: &v1.User{
			Id:        createdUser.ID.Value().String(),
			Name:      createdUser.Name.Value(),
			Email:     createdUser.Email.Value(),
			Avatar:    createdUser.Avatar.Value(),
			CreatedAt: timestamppb.New(createdUser.CreatedAt),
			UpdatedAt: timestamppb.New(createdUser.UpdatedAt),
		},
	}, nil
}

func (r *UserService) AllUsers(ctx context.Context, _ *v1.AllUsersRequest) (*v1.AllUsersResponse, error) {
	allUsers, err := r.usecase.All(ctx)
	if err != nil {
		r.logger.Warn("grpc - v1 - AllUsers", zap.Any("error:", err.Error()))

		return nil, fmt.Errorf("grpc - v1 - AllUsers: %w", err)
	}

	users := make([]*v1.User, len(allUsers))

	for i, u := range allUsers {
		users[i] = &v1.User{
			Id:        u.ID.Value().String(),
			Name:      u.Name.Value(),
			Email:     u.Email.Value(),
			Avatar:    u.Avatar.Value(),
			CreatedAt: timestamppb.New(u.CreatedAt),
			UpdatedAt: timestamppb.New(u.UpdatedAt),
		}
	}

	return &v1.AllUsersResponse{
		Users: users,
	}, nil
}
