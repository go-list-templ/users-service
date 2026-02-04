package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	redisrepo "github.com/go-list-templ/grpc/internal/adapter/cache/redis/repo"
	grpcserver "github.com/go-list-templ/grpc/internal/adapter/grpc/server"
	grpchandler "github.com/go-list-templ/grpc/internal/adapter/grpc/server/handler"
	httpserver "github.com/go-list-templ/grpc/internal/adapter/http/server"
	httphandler "github.com/go-list-templ/grpc/internal/adapter/http/server/handler"
	pgrepo "github.com/go-list-templ/grpc/internal/adapter/persistence/postgres/repo"

	"github.com/go-list-templ/grpc/internal/adapter/cache/redis"
	"github.com/go-list-templ/grpc/internal/adapter/persistence/postgres"
	"github.com/go-list-templ/grpc/internal/adapter/persistence/postgres/transaction"
	"github.com/go-list-templ/grpc/internal/core/service"
	"github.com/go-list-templ/grpc/pkg/config"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Panic(err)
	}
}

// nolint:errcheck,gocyclo,cyclop
func run() error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("starting app")
	logger.Info("initializing config")

	cfg, err := config.Load()
	if err != nil {
		logger.Panic("cant init config", zap.Error(err))
	}

	logger.Info("named logger")

	logger = logger.Named(cfg.App.Name)

	logger.Info("initializing postgres")

	pg, err := postgres.New(&cfg.DB, logger.With(zap.String("module", "postgres")))
	if err != nil {
		logger.Panic("cant init postgres", zap.Error(err))
	}
	defer pg.Close()

	logger.Info("postgres migration")

	migration := postgres.NewMigration(pg, logger.With(zap.String("module", "migration")))
	if err = migration.Up(); err != nil {
		logger.Panic("migration failed", zap.Error(err))
	}

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

	logger.Info("initializing transaction manager")

	trManager := transaction.NewManager(pg, logger.With(zap.String("module", "trx")))
	trGetter := transaction.NewTrmGetter(trManager)

	logger.Info("initializing repositories")

	outboxPostgresRepo := pgrepo.NewOutboxRepo(pg, trGetter)
	userPostgresRepo := pgrepo.NewUserRepo(pg, logger.With(zap.String("module", "pg user repo")), trGetter)
	userRedisRepo := redisrepo.NewUserRepo(userPostgresRepo, rd, logger.With(zap.String("module", "redis user repo")))

	logger.Info("initializing service")

	userService := service.NewUser(userRedisRepo, outboxPostgresRepo, trManager)

	logger.Info("initializing servers")

	grpcServer := grpcserver.New(&cfg.Server)
	grpcServer.Start()

	httpServer := httpserver.NewHTTP(&cfg.Server)
	httpServer.Start()

	logger.Info("registering handlers")

	grpchandler.RegisterUser(grpcServer.Server, userService, logger.With(zap.String("module", "user handler")))
	httphandler.RegisterDiagnostic(pg, rd, logger.With(zap.String("module", "diagnostic handler")))

	logger.Info("server started successfully")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case x := <-interrupt:
		logger.Info("Received a signal.", zap.String("signal", x.String()))
	case err = <-httpServer.Notify():
		logger.Error("Received an error from the health server", zap.Error(err))
	case err = <-grpcServer.Notify():
		logger.Error("Received an error from the grpc server", zap.Error(err))
	}

	logger.Info("stopping servers")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	grpcServer.Stop()

	if err = httpServer.Stop(ctx); err != nil {
		logger.Error("server stopped with error", zap.Error(err))
	}

	logger.Info("The app is calling the last defers and will be stopped")

	return nil
}
