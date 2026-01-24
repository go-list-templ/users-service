package user

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/domain/event"
	"github.com/go-list-templ/grpc/internal/repo"
	"github.com/go-list-templ/grpc/pkg/trm"
)

type Usecase struct {
	userRepo   repo.UserRepo
	outboxRepo repo.OutboxRepo

	trm *trm.Manager
}

func NewUserUsecase(u repo.UserRepo, o repo.OutboxRepo, t *trm.Manager) *Usecase {
	return &Usecase{u, o, t}
}

func (u *Usecase) Create(ctx context.Context, user entity.User) (entity.User, error) {
	err := u.trm.Do(ctx, func(ctx context.Context) error {
		err := u.userRepo.Store(ctx, user)
		if err != nil {
			return err
		}

		userCreated, err := event.NewUserCreated(user)
		if err != nil {
			return err
		}

		return u.outboxRepo.Publish(ctx, userCreated.Event)
	})
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *Usecase) All(ctx context.Context) ([]entity.User, error) {
	users, err := u.userRepo.All(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
