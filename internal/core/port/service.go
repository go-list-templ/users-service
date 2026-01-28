package port

import (
	"context"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
)

//go:generate mockgen -source=service.go -destination=../../../test/mocks/mock_service.go -package=mocks

type (
	UserService interface {
		Create(context.Context, entity.User) (entity.User, error)
		All(context.Context) ([]entity.User, error)
	}
)
