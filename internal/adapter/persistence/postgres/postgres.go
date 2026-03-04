package postgres

import (
	"context"
	"fmt"
	"github.com/exaring/otelpgx"
	"time"

	"github.com/go-list-templ/grpc/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	DefaultConnAttempts   = 10
	DefaultConnTimeout    = time.Second
	DefaultContextTimeout = 5 * time.Second
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
	conf.MaxConnIdleTime = cfg.MaxIdleTime

	connAttempts := DefaultConnAttempts
	connTimeout := DefaultConnTimeout

	conf.ConnConfig.Tracer = otelpgx.NewTracer()

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	for connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(ctx, conf)
		if err != nil {
			logger.Info("postgres err config", zap.Error(err))
		}

		err = pg.Ping(ctx)
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

	if err = otelpgx.RecordStats(pg.Pool); err != nil {
		return nil, fmt.Errorf("unable to record database stats: %w", err)
	}

	return pg, nil
}
