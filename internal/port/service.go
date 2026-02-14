package port

import (
	"context"
	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/go-list-templ/grpc/internal/core/dto"
)

//go:generate mockgen -source=service.go -destination=mock/mock_service.go -package=mock

type (
	UserService interface {
		Create(context.Context, dto.UserCreateInput) (entity.User, error)
		List(context.Context, dto.UserListInput) ([]entity.User, error)
	}
)
