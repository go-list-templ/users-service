package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/go-list-templ/users-service/pkg/config"
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

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	conf, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	conf.MaxConns = cfg.MaxConn
	conf.MaxConnLifetime = cfg.MaxConnTime
	conf.MaxConnIdleTime = cfg.MaxIdleTime
	conf.MinConns = cfg.MinConn
	conf.HealthCheckPeriod = cfg.HealthCheckTime
	conf.ConnConfig.ConnectTimeout = cfg.ConnTime
	conf.ConnConfig.Tracer = otelpgx.NewTracer()

	connAttempts := DefaultConnAttempts
	connTimeout := DefaultConnTimeout

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	for connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(ctx, conf)
		if err != nil {
			logger.Warn("postgres err config", zap.Error(err))
		}

		err = pg.Ping(ctx)
		if err == nil {
			break
		}

		logger.Warn("postgres is trying to connect", zap.Int("attempts", connAttempts), zap.Error(err))

		time.Sleep(connTimeout)

		connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("end attempts exceeded: %w", err)
	}

	return pg, nil
}
