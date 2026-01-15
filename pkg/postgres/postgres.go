package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/go-list-templ/grpc/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	DefaultConnAttempts = 10
	DefaultConnTimeout  = time.Second
)

type Postgres struct {
	*pgxpool.Pool
}

func New(cfg *config.DB, logger *zap.Logger) (*Postgres, error) {
	pg := &Postgres{}

	conf, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, err
	}

	conf.MaxConns = cfg.MaxConn
	conf.MaxConnLifetime = cfg.MaxConnTime

	connAttempts := DefaultConnAttempts
	connTimeout := DefaultConnTimeout

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for connAttempts > 0 {
		pg.Pool, _ = pgxpool.NewWithConfig(ctx, conf)

		err = pg.Pool.Ping(ctx)
		if err == nil {
			break
		}

		logger.Warn("Postgres is trying to connect", zap.Int("attempts", connAttempts))

		time.Sleep(connTimeout)

		connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("end attempts exceeded: %w", err)
	}

	return pg, nil
}
