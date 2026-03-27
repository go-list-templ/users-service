package repo

import (
	"context"
	"errors"
	"time"

	"github.com/go-list-templ/users-service/internal/adapter/cache/redis"
	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/dto"
	"github.com/go-list-templ/users-service/internal/port"
	"github.com/go-list-templ/users-service/pkg/paginate"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var ErrTypedSingleflight = errors.New("invalid type from singleflight")

const (
	TTLAllUsers = 10 * time.Minute

	TagAllUsers = "allUsers"
)

type User struct {
	repo   port.UserRepo
	redis  *redis.Redis
	logger *zap.Logger
	sf     singleflight.Group
}

func NewUser(repo port.UserRepo, redis *redis.Redis, logger *zap.Logger) *User {
	return &User{repo: repo, redis: redis, logger: logger, sf: singleflight.Group{}}
}

func (u *User) All(ctx context.Context, paginate paginate.Paginate) (dto.ListOutput, error) {
	var cached dto.ListOutput

	cacheKey := paginate.Cursor()
	pageToken := paginate.Token()

	if err := u.redis.GetCache(ctx, cacheKey, &cached); err == nil {
		u.logger.Info(
			"get from cache",
			zap.Any("context", ctx),
			zap.Any("page token", pageToken),
		)

		return cached, nil
	} else if !u.redis.ErrIsNil(err) {
		u.logger.Error(
			"get from cache",
			zap.Any("context", ctx),
			zap.Any("page token", pageToken),
			zap.Error(err),
		)
	}

	v, err, _ := u.sf.Do(cacheKey, func() (interface{}, error) {
		u.logger.Info(
			"get from persistent",
			zap.Any("context", ctx),
			zap.Any("page token", pageToken),
		)

		users, err := u.repo.All(ctx, paginate)
		if err != nil {
			return nil, err
		}

		if err = u.redis.SetByTags(ctx, cacheKey, users, TTLAllUsers, TagAllUsers); err != nil {
			u.logger.Warn(
				"set by tag",
				zap.Any("context", ctx),
				zap.Any("page token", pageToken),
				zap.Error(err),
			)
		}

		return users, nil
	})
	if err != nil {
		return dto.ListOutput{}, err
	}

	users, ok := v.(dto.ListOutput)
	if !ok {
		u.logger.Error(
			"singleflight typed",
			zap.Any("context", ctx),
			zap.Any("page token", pageToken),
			zap.Error(err),
		)

		return dto.ListOutput{}, ErrTypedSingleflight
	}

	return users, nil
}

func (u *User) Store(ctx context.Context, user entity.User) error {
	err := u.repo.Store(ctx, user)
	if err != nil {
		return err
	}

	if err = u.redis.InvalidateTags(ctx, TagAllUsers); err != nil {
		u.logger.Error("invalidate tag", zap.Any("context", ctx), zap.Error(err))
	}

	return nil
}
