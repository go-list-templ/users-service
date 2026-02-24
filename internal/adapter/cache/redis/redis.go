package redis

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/go-list-templ/grpc/pkg/config"
	"github.com/redis/go-redis/v9"
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
	})

	ctx, cancel := context.WithTimeout(context.Background(), DefaultCtx)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{client}, nil
}

func (r *Redis) DeleteCache(ctx context.Context, keys ...string) error {
	return r.Del(ctx, keys...).Err()
}

func (r *Redis) GetCache(ctx context.Context, key string, pointer any) error {
	data, err := r.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, pointer)
}

func (r *Redis) SetCache(ctx context.Context, key string, data any, ttl time.Duration) error {
	ttl = r.generateJitter(ttl)

	data, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) SetByTags(ctx context.Context, key string, data any, ttl time.Duration, tags ...string) error {
	pipe := r.TxPipeline()
	ttl = r.generateJitter(ttl)

	for _, tag := range tags {
		pipe.SAdd(ctx, tag, key)
		pipe.Expire(ctx, tag, ttl)
	}

	data, err := json.Marshal(data)
	if err != nil {
		return err
	}

	pipe.Set(ctx, key, data, ttl)

	_, errExec := pipe.Exec(ctx)
	return errExec
}

func (r *Redis) InvalidateTags(ctx context.Context, tags ...string) error {
	keys := make([]string, 0)

	for _, tag := range tags {
		k, _ := r.SMembers(ctx, tag).Result()
		keys = append(keys, tag)
		keys = append(keys, k...)
	}

	return r.DeleteCache(ctx, keys...)
}

func (r *Redis) generateJitter(ttl time.Duration) time.Duration {
	randomMultiplier := JitterMinFactor + rand.Float64()*(JitterMaxFactor-JitterMinFactor)

	jitter := time.Duration(float64(ttl) * randomMultiplier)

	return jitter
}
