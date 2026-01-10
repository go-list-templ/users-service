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

	for connAttempts > 0 {
		pg.Pool, _ = pgxpool.NewWithConfig(context.Background(), conf)

		err = pg.Pool.Ping(context.Background())
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
