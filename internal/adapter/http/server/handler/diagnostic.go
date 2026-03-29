package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-list-templ/users-service/internal/adapter/cache/redis"
	"github.com/go-list-templ/users-service/internal/adapter/persistence/postgres"
	"go.uber.org/zap"
)

const (
	TTL              = 30 * time.Second
	DefaultCtxTime   = 5 * time.Second
	MessageServerErr = "server error"
)

type Diagnostic struct {
	postgres *postgres.Postgres
	redis    *redis.Redis
	logger   *zap.Logger
}

func RegisterDiagnostic(postgres *postgres.Postgres, redis *redis.Redis, l *zap.Logger) {
	d := &Diagnostic{postgres, redis, l}

	http.HandleFunc("/healthz", d.HealthZ())
}

func (d *Diagnostic) HealthZ() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		var status int

		cacheKey := "healthz"

		ctx, cancel := context.WithTimeout(context.Background(), DefaultCtxTime)
		defer cancel()

		err := d.redis.GetCache(ctx, cacheKey, &status)
		if err == nil {
			data := map[string]int{
				"status": status,
			}

			d.writeJSON(w, status, data)

			return
		}

		status = http.StatusOK

		err = d.postgres.Ping(ctx)
		if err != nil {
			status = http.StatusServiceUnavailable

			d.logger.Warn("pinging postgres", zap.Error(err))
		}
		_, err = d.redis.Ping(ctx).Result()
		if err != nil {
			status = http.StatusServiceUnavailable

			d.logger.Warn("pinging redis", zap.Error(err))
		}

		err = d.redis.SetCache(ctx, cacheKey, status, TTL)
		if err != nil {
			d.logger.Warn("set cache", zap.Error(err))
		}

		data := map[string]int{
			"status": status,
		}

		d.writeJSON(w, status, data)
	}
}

func (d *Diagnostic) writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data, err := json.Marshal(v)
	if err != nil {
		d.logger.Error("json marshal", zap.Error(err))
		http.Error(w, MessageServerErr, http.StatusInternalServerError)

		return
	}

	_, err = w.Write(data)
	if err != nil {
		d.logger.Error("write json", zap.Error(err))
		http.Error(w, MessageServerErr, http.StatusInternalServerError)
	}
}
