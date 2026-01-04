package cache

import (
	"context"
	"encoding/json"
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
	logger zap.Logger
}

func NewUserRedisRepo(repo repo.UserRepo, redis redis.Redis, logger zap.Logger) *UserRedisRepo {
	return &UserRedisRepo{repo: repo, redis: redis, logger: logger}
}

func (u *UserRedisRepo) All(ctx context.Context) ([]entity.User, error) {
	var cachedUsers []dao.User

	cachedData, err := u.redis.Get(ctx, KeyAllUsers).Bytes()
	if err != nil {
		u.logger.Warn("failed to get cached users", zap.Error(err))
	}

	if err := json.Unmarshal(cachedData, &cachedUsers); err == nil {
		users := make([]entity.User, 0, len(cachedUsers))

		for _, user := range cachedUsers {
			users = append(users, user.ToEntity())
		}

		return users, nil
	}

	users, err := u.repo.All(ctx)
	if err != nil {
		return nil, err
	}

	go u.cacheAllUsers(users)
	u.logger.Info("all users from postgres")

	return users, nil
}

func (u *UserRedisRepo) cacheAllUsers(users []entity.User) {
	cacheUsers := make([]dao.User, 0, len(users))

	for _, user := range users {
		cacheUsers = append(cacheUsers, dao.FromEntity(user))
	}

	data, err := json.Marshal(cacheUsers)
	if err != nil {
		u.logger.Error("marshal users error", zap.Error(err))
	}

	err = u.redis.Set(context.Background(), KeyAllUsers, data, time.Hour).Err()
	if err != nil {
		u.logger.Error("redis set error", zap.Error(err))
	}
}

func (u *UserRedisRepo) Store(ctx context.Context, user entity.User) error {
	err := u.repo.Store(ctx, user)
	if err != nil {
		return err
	}

	go func() {
		if err = u.redis.Invalidate(KeyAllUsers); err != nil {
			u.logger.Error("redis del error", zap.Error(err))
		}
	}()

	return nil
}
