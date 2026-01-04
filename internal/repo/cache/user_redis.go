package cache

import (
	"context"
	"time"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/repo"
	"github.com/go-list-templ/grpc/internal/repo/cache/dao"
	"github.com/go-list-templ/grpc/pkg/redis"
	"go.uber.org/zap"
)

const (
	KeyAllUsers = "users:all"
)

type UserRedisRepo struct {
	repo   repo.UserRepo
	redis  redis.Redis
	logger *zap.Logger
}

func NewUserRedisRepo(repo repo.UserRepo, redis redis.Redis, logger zap.Logger) *UserRedisRepo {
	return &UserRedisRepo{repo: repo, redis: redis, logger: &logger}
}

func (u *UserRedisRepo) All(ctx context.Context) ([]entity.User, error) {
	var cachedUsers []dao.User

	err := u.redis.GetCache(KeyAllUsers, &cachedUsers)
	if err == nil && len(cachedUsers) > 0 {
		users := make([]entity.User, len(cachedUsers))

		for i, user := range cachedUsers {
			users[i] = user.ToEntity()
		}

		u.logger.Info("all users from cache")

		return users, nil
	}

	users, err := u.repo.All(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheUsers := make([]dao.User, len(users))

		for i, user := range users {
			cacheUsers[i] = dao.FromEntity(user)
		}

		err = u.redis.SetCache(KeyAllUsers, cacheUsers, time.Hour)
		if err != nil {
			u.logger.Error("redis set error", zap.Error(err))
		}
	}()

	u.logger.Info("all users from postgres")

	return users, nil
}

func (u *UserRedisRepo) Store(ctx context.Context, user entity.User) error {
	err := u.repo.Store(ctx, user)
	if err != nil {
		return err
	}

	go func() {
		if err = u.redis.DeleteCache(KeyAllUsers); err != nil {
			u.logger.Error("redis del error", zap.Error(err))
		}
	}()

	return nil
}
