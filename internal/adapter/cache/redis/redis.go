package redis

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"

	"github.com/go-list-templ/users-service/pkg/config"
	"github.com/klauspost/compress/s2"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	DefaultCtx = 5 * time.Second

	JitterMinFactor = 1.1
	JitterMaxFactor = 1.3
)

type Redis struct {
	*redis.Client
}

func New(cfg *config.Redis) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: "",
		DB:       0,

		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,

		MaxRetries:      cfg.MaxRetries,
		MinRetryBackoff: cfg.MinRetryBackoff,
		MaxRetryBackoff: cfg.MaxRetryBackoff,

		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleCons,
		PoolTimeout:  cfg.PoolTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), DefaultCtx)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	if err = redisotel.InstrumentTracing(client); err != nil {
		return nil, err
	}

	if err = redisotel.InstrumentMetrics(client); err != nil {
		return nil, err
	}

	return &Redis{client}, nil
}

func (r *Redis) ErrIsNil(err error) bool {
	return errors.Is(err, redis.Nil)
}

func (r *Redis) DeleteCache(ctx context.Context, keys ...string) error {
	return r.Del(ctx, keys...).Err()
}

func (r *Redis) GetCache(ctx context.Context, key string, pointer any) error {
	data, err := r.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	decompress, err := s2.Decode(nil, data)
	if err != nil {
		return err
	}

	return msgpack.Unmarshal(decompress, pointer)
}

func (r *Redis) SetCache(ctx context.Context, key string, data any, ttl time.Duration) error {
	ttl = r.generateJitter(ttl)

	pack, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}

	compress := s2.Encode(nil, pack)

	return r.Set(ctx, key, compress, ttl).Err()
}

func (r *Redis) SetNegativeCache(ctx context.Context, key string, ttl time.Duration) error {
	ttl = r.generateJitter(ttl)

	return r.Set(ctx, key, nil, ttl).Err()
}

func (r *Redis) SetByTags(ctx context.Context, key string, data any, ttl time.Duration, tags ...string) error {
	pipe := r.TxPipeline()
	ttl = r.generateJitter(ttl)

	for _, tag := range tags {
		pipe.SAdd(ctx, tag, key)
		pipe.Expire(ctx, tag, ttl)
	}

	pack, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}

	compress := s2.Encode(nil, pack)

	pipe.Set(ctx, key, compress, ttl)

	_, errExec := pipe.Exec(ctx)

	return errExec
}

func (r *Redis) InvalidateTags(ctx context.Context, tags ...string) error {
	keys := make([]string, 0)

	for _, tag := range tags {
		k, err := r.SMembers(ctx, tag).Result()
		if err != nil {
			return err
		}

		keys = append(keys, k...)
	}

	if len(keys) <= 0 {
		return nil
	}

	return r.DeleteCache(ctx, keys...)
}

func (r *Redis) generateJitter(ttl time.Duration) time.Duration {
	//nolint:gosec
	randomMultiplier := JitterMinFactor + rand.Float64()*(JitterMaxFactor-JitterMinFactor)

	jitter := time.Duration(float64(ttl) * randomMultiplier)

	return jitter
}
