package repo

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/entity"
)

type (
	UserStorageRepo interface {
		Store(context.Context, entity.User) error
		All(context.Context) ([]entity.User, error)
	}
)
