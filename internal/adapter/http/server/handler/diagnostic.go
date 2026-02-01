package handler

import (
	"context"
	"go.uber.org/zap"
	"net/http"

	"github.com/go-list-templ/grpc/internal/adapter/cache/redis"
	"github.com/go-list-templ/grpc/internal/adapter/persistence/postgres"
)

type Diagnostic struct {
	pg     *postgres.Postgres
	rd     *redis.Redis
	logger *zap.Logger
}

func RegisterDiagnostic(pg *postgres.Postgres, rd *redis.Redis, l *zap.Logger) {
	d := &Diagnostic{pg, rd, l}

	http.HandleFunc("/healthz", d.HealthZ())
}

func (d *Diagnostic) HealthZ() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK

		ctx := context.Background()

		err := d.pg.Ping(ctx)
		if err != nil {
			status = http.StatusServiceUnavailable

			d.logger.Warn("error pinging postgres", zap.Error(err))
		}

		_, err = d.rd.Ping(ctx).Result()
		if err != nil {
			status = http.StatusServiceUnavailable

			d.logger.Warn("error pinging redis", zap.Error(err))
		}

		d.logger.Info("status", zap.Int("status", status))

		w.WriteHeader(status)
	}
}
