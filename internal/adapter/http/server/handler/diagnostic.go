package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/go-list-templ/users-service/internal/adapter/cache/redis"
	"github.com/go-list-templ/users-service/internal/adapter/persistence/postgres"
	"go.uber.org/zap"
)

const DefaultCtxTime = 5 * time.Second

type Diagnostic struct {
	postgres *postgres.Postgres
	redis    *redis.Redis
	logger   *zap.Logger
}

func RegisterDiagnostic(postgres *postgres.Postgres, redis *redis.Redis, l *zap.Logger) {
	d := &Diagnostic{postgres, redis, l}

	http.HandleFunc("/healthz", d.Health())
	http.HandleFunc("/readyz", d.Ready())
}

func (d *Diagnostic) Health() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (d *Diagnostic) Ready() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), DefaultCtxTime)
		defer cancel()

		if err := d.postgres.Ping(ctx); err != nil {
			d.logger.Error("postgres unavailable", zap.Error(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		if err := d.redis.Ping(ctx).Err(); err != nil {
			d.logger.Error("redis unavailable", zap.Error(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
