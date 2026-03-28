package service

import (
	"context"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/domain/event"
	"github.com/go-list-templ/users-service/internal/core/domain/vo"
	"github.com/go-list-templ/users-service/internal/core/dto"
	"github.com/go-list-templ/users-service/internal/port"
	"github.com/go-list-templ/users-service/pkg/hasher"
	"github.com/go-list-templ/users-service/pkg/paginate"
)

type User struct {
	repo   port.UserRepo
	outbox port.OutboxRepo
	trm    port.TransactionManager
}

func NewUser(u port.UserRepo, o port.OutboxRepo, t port.TransactionManager) *User {
	return &User{u, o, t}
}

func (s *User) Create(ctx context.Context, input dto.CreateInput) (entity.User, error) {
	password, err := vo.NewPlainPassword(input.Password)
	if err != nil {
		return entity.User{}, err
	}

	hash, err := hasher.Hash(password.Value())
	if err != nil {
		return entity.User{}, err
	}

	user, err := entity.NewUser(input.Name, input.Email, hash)
	if err != nil {
		return entity.User{}, err
	}

	err = s.trm.Do(ctx, func(ctx context.Context) error {
		err = s.repo.Store(ctx, user)
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
	return s.repo.List(ctx, paginate.NewUUIDPaginate(input.PageToken))
}

func (s *User) GetByEmail(context.Context, dto.GetByEmailInput) (entity.User, error) {
	return entity.User{}, nil
}
