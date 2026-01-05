package user

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/repo"
)

type UseCase struct {
	repo       repo.UserRepo
	avatarRepo repo.UserAvatarRepo
}

func New(r repo.UserRepo, a repo.UserAvatarRepo) *UseCase {
	return &UseCase{repo: r, avatarRepo: a}
}

func (u *UseCase) All(ctx context.Context) ([]entity.User, error) {
	users, err := u.repo.All(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UseCase) Create(ctx context.Context, user entity.User) (entity.User, error) {
	err := u.repo.Store(ctx, user)
	if err != nil {
		return user, err
	}

	u.avatarRepo.Set(user)

	return user, nil
}
