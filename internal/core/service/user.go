package service

import (
	"context"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/go-list-templ/grpc/internal/core/domain/event"
	"github.com/go-list-templ/grpc/internal/core/dto"
	"github.com/go-list-templ/grpc/internal/port"
	"github.com/go-list-templ/grpc/pkg/paginate"
)

type User struct {
	userRepo   port.UserRepo
	outboxRepo port.OutboxRepo
	trm        port.TransactionManager
}

func NewUser(u port.UserRepo, o port.OutboxRepo, t port.TransactionManager) *User {
	return &User{u, o, t}
}

func (s *User) Create(ctx context.Context, input dto.UserCreateInput) (dto.UserCreateOutput, error) {
	user, err := entity.NewUser(input.Name, input.Email)
	if err != nil {
		return dto.UserCreateOutput{}, err
	}

	err = s.trm.Do(ctx, func(ctx context.Context) error {
		err = s.userRepo.Store(ctx, user)
		if err != nil {
			return err
		}

		userCreated, err := event.NewUserCreated(user)
		if err != nil {
			return err
		}

		return s.outboxRepo.Publish(ctx, userCreated.Event)
	})
	if err != nil {
		return dto.UserCreateOutput{}, err
	}

	return dto.UserCreateOutput{
		User: dto.UserFromEntity(user),
	}, nil
}

func (s *User) List(ctx context.Context, input dto.UserListInput) (dto.UserListOutput, error) {
	pagination := paginate.NewUUIDPaginate(input.PageSize, input.PageToken)

	users, err := s.userRepo.All(ctx, pagination)
	if err != nil {
		return dto.UserListOutput{}, err
	}

	pageToken := ""

	if len(users) > 0 {
		lastUser := users[len(users)-1]
		pageToken = pagination.GenerateToken(lastUser.ID.Value().String())
	}

	return dto.UserListOutput{
		Users:         dto.UsersFromEntity(users),
		NextPageToken: pageToken,
	}, nil
}
