package port

import (
	"context"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
)

type UserService interface {
	Create(context.Context, entity.User) (entity.User, error)
	All(context.Context) ([]entity.User, error)
}
