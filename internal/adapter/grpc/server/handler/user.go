package handler

import (
	"context"
	"fmt"

	v1 "github.com/go-list-templ/proto/gen/api/user/v1"
	pbgrpc "google.golang.org/grpc"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/go-list-templ/grpc/internal/core/port"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	v1.UserServiceServer

	userService port.UserService
	logger      *zap.Logger
}

func RegisterUser(s *pbgrpc.Server, u port.UserService, l *zap.Logger) {
	service := &User{userService: u, logger: l}
	{
		v1.RegisterUserServiceServer(s, service)
	}
}

func (u *User) CreateUser(ctx context.Context, request *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	user, err := entity.NewUser(request.GetName(), request.GetEmail())
	if err != nil {
		u.logger.Warn("CreateUser", zap.Error(err))
		return nil, fmt.Errorf("CreateUser: %w", err)
	}

	createdUser, err := u.userService.Create(ctx, user)
	if err != nil {
		u.logger.Warn("CreateUser", zap.Error(err))
		return nil, fmt.Errorf("CreateUser: %w", err)
	}

	return &v1.CreateUserResponse{
		User: u.toProto(createdUser),
	}, nil
}

func (u *User) AllUsers(ctx context.Context, _ *v1.AllUsersRequest) (*v1.AllUsersResponse, error) {
	allUsers, err := u.userService.All(ctx)
	if err != nil {
		u.logger.Warn("AllUsers", zap.Error(err))
		return nil, fmt.Errorf("AllUsers: %w", err)
	}

	users := make([]*v1.User, len(allUsers))

	for i, user := range allUsers {
		users[i] = u.toProto(user)
	}

	return &v1.AllUsersResponse{
		Users: users,
	}, nil
}

func (u *User) toProto(user entity.User) *v1.User {
	return &v1.User{
		Id:        user.ID.Value().String(),
		Name:      user.Name.Value(),
		Email:     user.Email.Value(),
		Avatar:    user.Avatar.Value(),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
