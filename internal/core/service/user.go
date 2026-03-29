package service

import (
	"context"
	"errors"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/domain/entityerr"
	"github.com/go-list-templ/users-service/internal/core/domain/event"
	"github.com/go-list-templ/users-service/internal/core/domain/vo"
	"github.com/go-list-templ/users-service/internal/core/dto"
	"github.com/go-list-templ/users-service/internal/port"
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
	user, err := entity.NewUser(input.Name, input.Email, input.Password)
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

func (s *User) GetByEmail(ctx context.Context, input dto.GetByEmailInput) (entity.User, error) {
	validEmail, err := vo.NewEmail(input.Email)
	if err != nil {
		return entity.User{}, entityerr.NewUserError("email", err)
	}

	user, err := s.repo.GetByEmail(ctx, validEmail)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *User) VerifyCred(ctx context.Context, input dto.VerifyCredInput) (entity.User, error) {
	validEmail, err := vo.NewEmail(input.Email)
	if err != nil {
		return entity.User{}, entityerr.NewUserError("email", err)
	}

	validPass, err := vo.NewPlainPassword(input.Password)
	if err != nil {
		return entity.User{}, entityerr.NewUserError("password", err)
	}

	user, err := s.repo.GetByEmail(ctx, validEmail)
	if errors.Is(err, entityerr.ErrUserNotFound) {
		return entity.User{}, entityerr.ErrUserVerifyCred
	}
	if err != nil {
		return entity.User{}, err
	}

	if !user.Password.Compare(validPass) {
		return entity.User{}, entityerr.ErrUserVerifyCred
	}

	return user, err
}
