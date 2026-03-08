package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/exaring/otelpgx"
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
	Master  *pgxpool.Pool
	Replica *pgxpool.Pool
}

func New(cfg *config.DB, logger *zap.Logger) (*Postgres, error) {
	master, err := initMaster(cfg, logger)
	if err != nil {
		logger.Error("init master", zap.Error(err))

		return nil, err
	}

	replica, err := initReplica(cfg, logger)
	if err != nil {
		logger.Error("init replica", zap.Error(err))

		return nil, err
	}

	return &Postgres{
		Master:  master,
		Replica: replica,
	}, nil
}

func (p *Postgres) Shutdown() {
	p.Master.Close()
	p.Replica.Close()
}

func initMaster(cfg *config.DB, logger *zap.Logger) (*pgxpool.Pool, error) {
	confMaster, err := pgxpool.ParseConfig(cfg.WriteURL)
	if err != nil {
		return nil, err
	}

	confMaster.MaxConns = cfg.MaxConn
	confMaster.MaxConnLifetime = cfg.MaxConnTime
	confMaster.MaxConnIdleTime = cfg.MaxIdleTime

	confMaster.ConnConfig.Tracer = otelpgx.NewTracer()

	return initPool(confMaster, logger)
}

func initReplica(cfg *config.DB, logger *zap.Logger) (*pgxpool.Pool, error) {
	confMaster, err := pgxpool.ParseConfig(cfg.ReadURL)
	if err != nil {
		return nil, err
	}

	confMaster.MaxConns = cfg.MaxConn
	confMaster.MaxConnLifetime = cfg.MaxConnTime
	confMaster.MaxConnIdleTime = cfg.MaxIdleTime

	confMaster.ConnConfig.Tracer = otelpgx.NewTracer()

	return initPool(confMaster, logger)
}

func initPool(conf *pgxpool.Config, logger *zap.Logger) (*pgxpool.Pool, error) {
	attempts := DefaultConnAttempts

	for attempts > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)

		pool, err := pgxpool.NewWithConfig(ctx, conf)
		if err == nil {
			err = pool.Ping(ctx)
		}
		cancel()

		if err == nil {
			if err = otelpgx.RecordStats(pool); err != nil {
				return nil, err
			}

			return pool, nil
		}

		logger.Warn("Postgres is trying to connect", zap.Int("attempts", attempts))

		attempts--
		time.Sleep(DefaultConnTimeout)
	}

	return nil, fmt.Errorf("all attempts failed")
}
