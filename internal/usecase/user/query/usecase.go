package query

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/repo"
)

type UserUsecase struct {
	userRepo repo.UserRepo
}

func NewUserUsecase(u repo.UserRepo) *UserUsecase {
	return &UserUsecase{userRepo: u}
}

func (u *UserUsecase) All(ctx context.Context) ([]entity.User, error) {
	users, err := u.userRepo.All(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
