package service

import (
	"context"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/domain/event"
	"github.com/go-list-templ/users-service/internal/core/dto"
	"github.com/go-list-templ/users-service/internal/port"
	"github.com/go-list-templ/users-service/pkg/paginate"
)

type User struct {
	user   port.UserRepo
	outbox port.OutboxRepo
	trm    port.TransactionManager
}

func NewUser(u port.UserRepo, o port.OutboxRepo, t port.TransactionManager) *User {
	return &User{u, o, t}
}

func (s *User) Create(ctx context.Context, input dto.CreateInput) (entity.User, error) {
	user, err := entity.NewUser(input.Name, input.Email)
	if err != nil {
		return entity.User{}, err
	}

	err = s.trm.Do(ctx, func(ctx context.Context) error {
		err = s.user.Store(ctx, user)
		if err != nil {
			return err
		}

		userCreated, err := event.NewUserCreated(user)
		if err != nil {
			return err
		}

		return s.outbox.Publish(ctx, userCreated.Event)
	})
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *User) List(ctx context.Context, input dto.ListInput) (dto.ListOutput, error) {
	pagination := paginate.NewUUIDPaginate(input.PageToken)

	return s.user.List(ctx, pagination)
}
