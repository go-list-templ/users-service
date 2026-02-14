package handler

import (
	"context"
	"errors"
	"github.com/go-list-templ/grpc/internal/core/dto"
	v1 "github.com/go-list-templ/proto/gen/api/user/v1"
	pbgrpc "google.golang.org/grpc"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
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
	inputDTO := dto.UserCreateInput{
		Name:  request.GetName(),
		Email: request.GetEmail(),
	}

	outputDTO, err := u.userService.Create(ctx, inputDTO)
	if err != nil {
		u.logger.Warn("user service create", zap.Error(err))

		return nil, u.toGRPCError(err)
	}

	return &v1.CreateResponse{
		User: u.toProto(outputDTO),
	}, nil
}

func (u *User) List(ctx context.Context, request *v1.ListRequest) (*v1.ListResponse, error) {
	inputDTO := dto.UserListInput{
		PageSize:  request.GetPageSize(),
		PageToken: request.GetPageToken(),
	}

	outputDTO, err := u.userService.List(ctx, inputDTO)
	if err != nil {
		u.logger.Warn("all user", zap.Error(err))

		return nil, u.toGRPCError(err)
	}

	users := make([]*v1.User, len(outputDTO))

	for i, user := range outputDTO {
		users[i] = u.toProto(user)
	}

	return &v1.ListResponse{
		Users: users,
	}, nil
}

func (u *User) toProto(user dto.User) *v1.User {
	return &v1.User{
		Id:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func (u *User) toGRPCError(err error) error {
	if errors.Is(err, entity.ErrUserAlreadyExists) {
		return status.Error(codes.AlreadyExists, err.Error())
	}

	if errors.Is(err, entity.ErrUserNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, entity.ErrUserInvalidData) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return status.Error(codes.Internal, "internal error")
}
