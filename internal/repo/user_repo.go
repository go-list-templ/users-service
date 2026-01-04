package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/repo/dao"
	"github.com/go-list-templ/grpc/pkg/redis"
	"go.uber.org/zap"
)

const (
	CacheKeyAllUsers = "users:all"
)

type UserRepo struct {
	persistent UserStorageRepo
	redis      redis.Redis
	logger     zap.Logger
}

func NewUserRepo(p UserStorageRepo, r redis.Redis, logger zap.Logger) *UserRepo {
	return &UserRepo{persistent: p, redis: r, logger: logger}
}

func (u *UserRepo) All(ctx context.Context) ([]entity.User, error) {
	var cachedUsers []dao.User

	cachedData, err := u.redis.Get(ctx, CacheKeyAllUsers).Bytes()
	if err != nil {
		u.logger.Warn("failed to get cached users", zap.Error(err))
	}

	if err := json.Unmarshal(cachedData, &cachedUsers); err == nil {
		users := make([]entity.User, 0, len(cachedUsers))

		var errToEntity error

		for _, user := range cachedUsers {
			toEntity, err := user.ToEntity()
			if err != nil {
				errToEntity = err

				u.logger.Warn("failed to convert to entity", zap.Error(err), zap.Any("user", user))

				break
			}

			users = append(users, *toEntity)
		}

		if errToEntity == nil {
			u.logger.Info("all users from cache")

			return users, nil
		}
	}

	users, err := u.persistent.All(ctx)
	if err != nil {
		return nil, err
	}

	go u.cacheAllUsers(users)
	u.logger.Info("all users from postgres")

	return users, nil
}

func (u *UserRepo) cacheAllUsers(users []entity.User) {
	newCachedUsers := make([]dao.User, 0, len(users))
	for _, user := range users {
		newCachedUsers = append(newCachedUsers, dao.User{
			ID:        user.ID.Value(),
			Name:      user.Name.Value(),
			Email:     user.Email.Value(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	data, err := json.Marshal(newCachedUsers)
	if err != nil {
		u.logger.Error("marshal users error", zap.Error(err))
	}

	err = u.redis.Set(context.Background(), CacheKeyAllUsers, data, time.Hour).Err()
	if err != nil {
		u.logger.Error("redis set error", zap.Error(err))
	}
}

func (u *UserRepo) Create(ctx context.Context, user entity.User) error {
	err := u.persistent.Store(ctx, user)
	if err != nil {
		return err
	}

	go func() {
		err = u.redis.Del(context.Background(), CacheKeyAllUsers).Err()
		if err != nil {
			u.logger.Error("redis del error", zap.Error(err))
		}
	}()

	return nil
}
