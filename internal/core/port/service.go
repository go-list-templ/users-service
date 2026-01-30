package port

import (
	"context"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
)

//go:generate mockgen -source=service.go -destination=mock/mock_service.go -package=mock

type (
	UserService interface {
		Create(context.Context, entity.User) (entity.User, error)
		All(context.Context) ([]entity.User, error)
	}
)
