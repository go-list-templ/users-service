package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-list-templ/users-service/internal/adapter/cache/redis"
	"github.com/go-list-templ/users-service/internal/adapter/cache/redis/repo/dao"
	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/domain/entityerr"
	"github.com/go-list-templ/users-service/internal/core/domain/vo"
	"github.com/go-list-templ/users-service/internal/core/dto"
	"github.com/go-list-templ/users-service/internal/port"
	"github.com/go-list-templ/users-service/pkg/hasher"
	"github.com/go-list-templ/users-service/pkg/paginate"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var (
	ErrTypedSingleflight = errors.New("invalid type from singleflight")
)

const (
	TTLList       = 10 * time.Minute
	TTLGetByEmail = 15 * time.Minute
	TTLNegative   = 5 * time.Minute

	TagList    = "user-list"
	TagByEmail = "user-email"

	byEmailPrefix = "user:email"
	listKeyPrefix = "user:list"
)

func byEmailKey(hash string) string {
	return fmt.Sprintf("%s:%s", byEmailPrefix, hash)
}

func listKey(userID string) string {
	return fmt.Sprintf("%s:%s", listKeyPrefix, userID)
}

type User struct {
	repo   port.UserRepo
	redis  *redis.Redis
	logger *zap.Logger
	sf     singleflight.Group
}

func NewUser(repo port.UserRepo, redis *redis.Redis, logger *zap.Logger) *User {
	return &User{repo: repo, redis: redis, logger: logger, sf: singleflight.Group{}}
}

func (u *User) List(ctx context.Context, paginate paginate.Paginate) (dto.ListOutput, error) {
	var cached dto.ListOutput

	cacheKey := listKey(paginate.Cursor())
	pageToken := paginate.Token()

	if err := u.redis.GetCache(ctx, cacheKey, &cached); err == nil {
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

		users, err := u.repo.List(ctx, paginate)
		if err != nil {
			return nil, err
		}

		if err = u.redis.SetByTags(ctx, cacheKey, users, TTLList, TagList); err != nil {
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

	keys := []string{
		byEmailKey(hasher.EmailHash(user.Email.Value())),
	}

	if err = u.redis.InvalidateTags(ctx, TagList); err != nil {
		u.logger.Error("invalidate tags", zap.Any("context", ctx), zap.Error(err))
	}

	if err = u.redis.DeleteCache(ctx, keys...); err != nil {
		u.logger.Error("delete cache", zap.Any("context", ctx), zap.Error(err))
	}

	return nil
}

func (u *User) GetByEmail(ctx context.Context, email vo.Email) (entity.User, error) {
	var cached dao.User

	cacheKey := byEmailKey(hasher.EmailHash(email.Value()))

	if err := u.redis.GetCache(ctx, cacheKey, &cached); err == nil {
		if cached.IsEmpty() {
			return entity.User{}, entityerr.ErrUserNotFound
		}

		return cached.ToEntity(), nil
	} else if !u.redis.ErrIsNil(err) {
		u.logger.Error(
			"get from cache",
			zap.Any("context", ctx),
			zap.Any("cache key", cacheKey),
			zap.Error(err),
		)
	}

	v, err, _ := u.sf.Do(cacheKey, func() (interface{}, error) {
		u.logger.Info(
			"get from persistent",
			zap.Any("context", ctx),
			zap.Any("cache key", cacheKey),
		)

		user, err := u.repo.GetByEmail(ctx, email)
		if errors.Is(err, entityerr.ErrUserNotFound) {
			if err = u.redis.SetCache(ctx, cacheKey, cached, TTLNegative); err != nil {
				u.logger.Warn(
					"set negative cache",
					zap.Any("context", ctx),
					zap.Any("cache key", cacheKey),
					zap.Error(err),
				)
			}

			return nil, entityerr.ErrUserNotFound
		}
		if err != nil {
			return nil, err
		}

		if err = u.redis.SetByTags(ctx, cacheKey, dao.FromEntity(user), TTLGetByEmail, TagByEmail); err != nil {
			u.logger.Warn(
				"set cache",
				zap.Any("context", ctx),
				zap.Any("cache key", cacheKey),
				zap.Error(err),
			)
		}

		return user, nil
	})
	if err != nil {
		return entity.User{}, err
	}

	user, ok := v.(entity.User)
	if !ok {
		u.logger.Error(
			"singleflight typed",
			zap.Any("context", ctx),
			zap.Any("cache key", cacheKey),
			zap.Error(err),
		)

		return entity.User{}, ErrTypedSingleflight
	}

	return user, nil
}
