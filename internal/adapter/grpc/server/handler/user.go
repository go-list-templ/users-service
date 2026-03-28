package handler

import (
	"context"

	v1 "github.com/go-list-templ/proto/gen/api/user/v1"
	pbgrpc "google.golang.org/grpc"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/dto"
	"github.com/go-list-templ/users-service/internal/port"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	v1.UserServiceServer

	service port.UserService
	logger  *zap.Logger
}

func RegisterUser(s *pbgrpc.Server, u port.UserService, l *zap.Logger) {
	service := &User{service: u, logger: l}
	{
		v1.RegisterUserServiceServer(s, service)
	}
}

func (u *User) Create(ctx context.Context, request *v1.CreateRequest) (*v1.CreateResponse, error) {
	input := dto.CreateInput{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	user, err := u.service.Create(ctx, input)
	if err != nil {
		u.logger.Warn("user create", zap.Any("context", ctx), zap.Error(err))

		return nil, err
	}

	return &v1.CreateResponse{
		User: u.entityToProto(user),
	}, nil
}

func (u *User) GetByEmail(ctx context.Context, request *v1.GetByEmailRequest) (*v1.GetByEmailResponse, error) {
	input := dto.GetByEmailInput{
		Email: request.GetEmail(),
	}

	user, err := u.service.GetByEmail(ctx, input)
	if err != nil {
		u.logger.Warn("user get by email", zap.Any("context", ctx), zap.Error(err))

		return nil, err
	}

	return &v1.GetByEmailResponse{
		User: u.entityToProto(user),
	}, nil
}

func (u *User) List(ctx context.Context, request *v1.ListRequest) (*v1.ListResponse, error) {
	input := dto.ListInput{
		PageToken: request.GetPageToken(),
	}

	output, err := u.service.List(ctx, input)
	if err != nil {
		u.logger.Warn("user list", zap.Any("context", ctx), zap.Error(err))

		return nil, err
	}

	users := make([]*v1.User, len(output.Users))

	for i, user := range output.Users {
		users[i] = u.dtoToProto(user)
	}

	return &v1.ListResponse{
		Users:         users,
		NextPageToken: output.NextPageToken,
	}, nil
}

func (u *User) dtoToProto(user dto.User) *v1.User {
	return &v1.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func (u *User) entityToProto(user entity.User) *v1.User {
	var nameVal *string
	if name, ok := user.Name.Get(); ok {
		str := name.Value()
		nameVal = &str
	}

	return &v1.User{
		Id:        user.ID.Value().String(),
		Name:      nameVal,
		Email:     user.Email.Value(),
		Avatar:    user.Avatar.Value(),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
