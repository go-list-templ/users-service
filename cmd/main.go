package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	grpcserver "github.com/go-list-templ/grpc/pkg/grpc/server"
	httpserver "github.com/go-list-templ/grpc/pkg/http/server"

	"github.com/go-list-templ/grpc/config"
	"github.com/go-list-templ/grpc/internal/controller/grpc"
	"github.com/go-list-templ/grpc/internal/repo/cache"
	"github.com/go-list-templ/grpc/internal/repo/storage"
	"github.com/go-list-templ/grpc/internal/usecase/user/command"
	"github.com/go-list-templ/grpc/internal/usecase/user/query"
	"github.com/go-list-templ/grpc/pkg/postgres"
	"github.com/go-list-templ/grpc/pkg/redis"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Panic(err)
	}
}

// nolint:errcheck
func run() error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("starting app")
	logger.Info("initializing config")

	cfg, err := config.Load()
	if err != nil {
		logger.Panic("cant init config", zap.Error(err))
	}

	logger.Info("initializing postgres")

	pg, err := postgres.New(&cfg.DB, logger)
	if err != nil {
		logger.Panic("cant init postgres", zap.Error(err))
	}
	defer pg.Close()

	logger.Info("initializing redis")

	rd, err := redis.New(&cfg.Redis)
	if err != nil {
		logger.Panic("cant init redis", zap.Error(err))
	}
	defer func() {
		if err = rd.Close(); err != nil {
			logger.Error("redis close failed", zap.Error(err))
		}
	}()

	logger.Info("initializing http client")

	// hc := httpclient.New(cfg.Client)

	logger.Info("initializing repositories")

	outboxPostgresRepo := storage.NewOutboxPostgres(pg)
	userPostgresRepo := storage.NewUserPostgres(pg)
	userRedisRepo := cache.NewUserRedis(userPostgresRepo, rd, logger)
	// userUnavatarRepo := external.NewUserUnavatar(hc, logger)

	logger.Info("initializing usecase")

	userQueryUsecase := query.NewUserUsecase(userRedisRepo)
	userCommandUsecase := command.NewUserUsecase(userRedisRepo, outboxPostgresRepo)

	logger.Info("initializing servers")

	grpcServer := grpcserver.NewAPI(&cfg.Server)
	grpcServer.Start()

	healthServer := httpserver.NewHealth(&cfg.Server)
	healthServer.Start()

	logger.Info("initializing routes")

	grpcRouter := grpc.NewRouter(grpcServer.Server, logger)
	grpcRouter.UserRegister(userQueryUsecase, userCommandUsecase)

	grpcRouter.Release()

	logger.Info("server started successfully")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case x := <-interrupt:
		logger.Info("Received a signal.", zap.String("signal", x.String()))
	case err = <-healthServer.Notify():
		logger.Error("Received an error from the health server", zap.Error(err))
	case err = <-grpcServer.Notify():
		logger.Error("Received an error from the grpc server", zap.Error(err))
	}

	logger.Info("stopping servers")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	grpcServer.Stop()

	if err = healthServer.Stop(ctx); err != nil {
		logger.Error("server stopped with error", zap.Error(err))
	}

	logger.Info("The app is calling the last defers and will be stopped")

	return nil
}
