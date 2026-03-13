package port

import (
	"context"
	"github.com/go-list-templ/users-service/internal/core/dto"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/domain/event"
	"github.com/go-list-templ/users-service/pkg/paginate"
)

//go:generate mockgen -source=repo.go -destination=mock/mock_repo.go -package=mock

type (
	UserRepo interface {
		Store(context.Context, entity.User) error
		All(context.Context, paginate.Paginate) (dto.UserListOutput, error)
	}

	OutboxRepo interface {
		Publish(context.Context, event.Event) error
	}
)
