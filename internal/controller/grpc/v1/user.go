package v1

import (
	"context"
	"fmt"

	userUsecase "github.com/go-list-templ/grpc/internal/usecase/user"
	v1 "github.com/go-list-templ/proto/gen/api/user/v1"
	pbgrpc "google.golang.org/grpc"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	v1.UserServiceServer

	usecase *userUsecase.Usecase
	logger  *zap.Logger
}

func NewUserRoute(server *pbgrpc.Server, u *userUsecase.Usecase, l *zap.Logger) {
	service := &UserService{usecase: u, logger: l}
	{
		v1.RegisterUserServiceServer(server, service)
	}
}

func (u *UserService) CreateUser(ctx context.Context, request *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	user, err := entity.NewUser(request.GetName(), request.GetEmail())
	if err != nil {
		u.logger.Warn("grpc - v1 - NewUser", zap.Any("error:", err.Error()))

		return nil, fmt.Errorf("grpc - v1 - NewUser: %w", err)
	}

	createdUser, err := u.usecase.Create(ctx, *user)
	if err != nil {
		u.logger.Warn("grpc - v1 - CreateUser", zap.Any("error:", err.Error()))

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

func (u *UserService) AllUsers(ctx context.Context, _ *v1.AllUsersRequest) (*v1.AllUsersResponse, error) {
	allUsers, err := u.usecase.All(ctx)
	if err != nil {
		u.logger.Warn("grpc - v1 - AllUsers", zap.Any("error:", err.Error()))

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
