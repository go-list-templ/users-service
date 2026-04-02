package main

import (
	"context"
	"github.com/go-list-templ/users-service/pkg/migrator"
	"log"
	"os"
	"os/signal"
	"syscall"

	redisrepo "github.com/go-list-templ/users-service/internal/adapter/cache/redis/repo"
	grpcserver "github.com/go-list-templ/users-service/internal/adapter/grpc/server"
	grpchandler "github.com/go-list-templ/users-service/internal/adapter/grpc/server/handler"
	httpserver "github.com/go-list-templ/users-service/internal/adapter/http/server"
	httphandler "github.com/go-list-templ/users-service/internal/adapter/http/server/handler"
	pgrepo "github.com/go-list-templ/users-service/internal/adapter/persistence/postgres/repo"

	"github.com/go-list-templ/users-service/internal/adapter/cache/redis"
	"github.com/go-list-templ/users-service/internal/adapter/persistence/postgres"
	"github.com/go-list-templ/users-service/internal/adapter/persistence/postgres/transaction"
	"github.com/go-list-templ/users-service/internal/core/service"
	"github.com/go-list-templ/users-service/pkg/config"
	"github.com/go-list-templ/users-service/pkg/otel"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Panic(err)
	}
}

// nolint:errcheck,gocyclo,cyclop
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	telemetry, err := otel.NewTelemetry(cfg)
	if err != nil {
		return err
	}

	logger := telemetry.Logger.Logger

	logger.Info("starting app",
		zap.String("name", cfg.App.Name),
		zap.String("version", cfg.App.Version),
	)

	maxProcsShowdown, err := maxprocs.Set(maxprocs.Logger(func(_ string, args ...interface{}) {
		logger.Info("auto max procs", zap.Any("count", args))
	}))
	if err != nil {
		logger.Error("set auto max procs", zap.Error(err))
	}

	logger.Info("initializing postgres")

	pg, err := postgres.New(&cfg.DB, logger.With(zap.String("module", "postgres")))
	if err != nil {
		logger.Panic("init postgres", zap.Error(err))
	}

	logger.Info("migrations up")

	if err = migrator.Up(&cfg.DB); err != nil {
		logger.Panic("migrations up", zap.Error(err))
	}

	logger.Info("initializing redis")

	rd, err := redis.New(&cfg.Redis)
	if err != nil {
		logger.Panic("init redis", zap.Error(err))
	}

	logger.Info("initializing transaction manager")

	trManager := transaction.NewManager(pg, logger.With(zap.String("module", "trx")))
	trGetter := transaction.NewTrmGetter(trManager)

	logger.Info("initializing repositories")

	outboxPostgresRepo := pgrepo.NewOutboxRepo(pg, trGetter)
	userPostgresRepo := pgrepo.NewUser(pg, logger.With(zap.String("module", "pg user repo")), trGetter)
	userRedisRepo := redisrepo.NewUser(userPostgresRepo, rd, logger.With(zap.String("module", "redis user repo")))

	logger.Info("initializing services")

	userService := service.NewUser(userRedisRepo, outboxPostgresRepo, trManager)

	logger.Info("initializing servers")

	grpcServer := grpcserver.New(&cfg.Server, logger.With(zap.String("module", "grpc server")))
	grpcServer.Start()

	httpServer := httpserver.NewHTTP(&cfg.Server)
	httpServer.Start()

	logger.Info("registering grpc handlers")

	grpchandler.RegisterUser(grpcServer.Server, userService, logger.With(zap.String("module", "user handler")))

	logger.Info("registering http handlers")

	httphandler.RegisterDiagnostic(pg, rd, logger.With(zap.String("module", "diagnostic handler")))

	logger.Info("server started successfully")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case x := <-interrupt:
		logger.Info("Received a signal.", zap.String("signal", x.String()))
	case err = <-httpServer.Notify():
		logger.Error("Received from the http server", zap.Error(err))
	case err = <-grpcServer.Notify():
		logger.Error("Received from the grpc server", zap.Error(err))
	}

	logger.Info("stopping app")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err = grpcServer.Shutdown(); err != nil {
		logger.Error("grpc shutdown", zap.Error(err))
	}

	if err = httpServer.Shutdown(ctx); err != nil {
		logger.Error("http shutdown", zap.Error(err))
	}

	if err = rd.Close(); err != nil {
		logger.Error("redis close", zap.Error(err))
	}

	pg.Close()

	if err = telemetry.Shutdown(ctx); err != nil {
		logger.Error("telemetry shutdown", zap.Error(err))
	}

	maxProcsShowdown()

	return nil
}
