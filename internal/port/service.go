package port

import (
	"context"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/dto"
)

//go:generate mockgen -source=service.go -destination=mock/mock_service.go -package=mock

type (
	UserService interface {
		Create(context.Context, dto.UserCreateInput) (entity.User, error)
		List(context.Context, dto.UserListInput) (dto.UserListOutput, error)
	}
)
