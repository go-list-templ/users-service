package handler

import (
	"context"
	"errors"

	v1 "github.com/go-list-templ/proto/gen/api/user/v1"
	pbgrpc "google.golang.org/grpc"

	"github.com/go-list-templ/grpc/internal/core/domain/entityerr"
	"github.com/go-list-templ/grpc/internal/core/dto"
	"github.com/go-list-templ/grpc/internal/port"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (u *User) Create(ctx context.Context, request *v1.CreateRequest) (*v1.CreateResponse, error) {
	input := dto.UserCreateInput{
		Name:  request.GetName(),
		Email: request.GetEmail(),
	}

	output, err := u.userService.Create(ctx, input)
	if err != nil {
		u.logger.Warn("user service create", zap.Error(err))

		return nil, u.toGRPCError(err)
	}

	return u.createToProto(output), nil
}

func (u *User) List(ctx context.Context, request *v1.ListRequest) (*v1.ListResponse, error) {
	input := dto.UserListInput{
		PageToken: request.GetPageToken(),
	}

	output, err := u.userService.List(ctx, input)
	if err != nil {
		u.logger.Warn("all user", zap.Error(err))

		return nil, u.toGRPCError(err)
	}

	return u.listToProto(output), nil
}

func (u *User) userToProto(user dto.User) *v1.User {
	return &v1.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func (u *User) createToProto(output dto.UserCreateOutput) *v1.CreateResponse {
	return &v1.CreateResponse{
		User: u.userToProto(output.User),
	}
}

func (u *User) listToProto(output dto.UserListOutput) *v1.ListResponse {
	users := make([]*v1.User, len(output.Users))

	for i, user := range output.Users {
		users[i] = u.userToProto(user)
	}

	return &v1.ListResponse{
		Users:         users,
		NextPageToken: output.NextPageToken,
	}
}

func (u *User) toGRPCError(err error) error {
	switch {
	case errors.Is(err, entityerr.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, entityerr.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, entityerr.ErrUserInvalidData):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
