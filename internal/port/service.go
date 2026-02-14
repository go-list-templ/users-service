package port

import (
	"context"
	"github.com/go-list-templ/grpc/internal/core/dto"
)

//go:generate mockgen -source=service.go -destination=mock/mock_service.go -package=mock

type (
	UserService interface {
		Create(context.Context, dto.UserCreateInput) (dto.User, error)
		List(context.Context, dto.UserListInput) ([]dto.User, error)
	}
)
