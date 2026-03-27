package port

import (
	"context"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/domain/event"
	"github.com/go-list-templ/users-service/internal/core/dto"
	"github.com/go-list-templ/users-service/pkg/paginate"
)

//go:generate mockgen -source=repo.go -destination=mock/mock_repo.go -package=mock

type (
	UserRepo interface {
		GetByEmail(context.Context, dto.GetByEmailInput) (entity.User, error)
		List(context.Context, paginate.Paginate) (dto.ListOutput, error)
		Store(context.Context, entity.User) error
	}

	OutboxRepo interface {
		Publish(context.Context, event.Event) error
	}
)
