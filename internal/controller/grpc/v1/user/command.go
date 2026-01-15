package user

import (
	"context"
	"fmt"

	v1 "github.com/go-list-templ/proto/gen/api/user/v1"
	pbgrpc "google.golang.org/grpc"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/usecase/user/command"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CommandService struct {
	v1.UserCommandServiceServer

	usecase *command.UserUsecase
	logger  *zap.Logger
}

func NewCommandService(server *pbgrpc.Server, u *command.UserUsecase, l *zap.Logger) {
	service := &CommandService{usecase: u, logger: l}
	{
		v1.RegisterUserCommandServiceServer(server, service)
	}
}

func (c *CommandService) CreateUser(ctx context.Context, request *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	user, err := entity.NewUser(request.GetName(), request.GetEmail())
	if err != nil {
		c.logger.Warn("grpc - v1 - NewUser", zap.Any("error:", err.Error()))

		return nil, fmt.Errorf("grpc - v1 - NewUser: %w", err)
	}

	createdUser, err := c.usecase.Create(ctx, *user)
	if err != nil {
		c.logger.Warn("grpc - v1 - CreateUser", zap.Any("error:", err.Error()))

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
