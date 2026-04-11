package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	App struct {
		Name    string `envconfig:"APP_NAME"`
		Version string `envconfig:"APP_VERSION"`
	}

	Server struct {
		GRPCPort        string        `envconfig:"GRPC_PORT"`
		GRPCTime        time.Duration `envconfig:"GRPC_TIME"`
		GRPCTimeout     time.Duration `envconfig:"GRPC_TIMEOUT"`
		GRPCMaxConnIdle time.Duration `envconfig:"GRPC_MAX_CONN_IDLE"`
		GRPCMaxConnAge  time.Duration `envconfig:"GRPC_MAX_CONN_AGE"`

		HTTPort     string        `envconfig:"HTTP_PORT"`
		HTTPTimeout time.Duration `envconfig:"HTTP_TIMEOUT"`
		IdleTimeout time.Duration `envconfig:"IDLE_TIMEOUT"`

		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT"`
	}

	DB struct {
		Host     string `envconfig:"DB_HOST"`
		Port     uint16 `envconfig:"DB_PORT"`
		Name     string `envconfig:"DB_NAME"`
		Username string `envconfig:"DB_USERNAME"`
		Password string `envconfig:"DB_PASSWORD"`

		MaxConn     int32         `envconfig:"DB_MAX_CONN"`
		MinConn     int32         `envconfig:"DB_MIN_CONN"`
		ConnTime    time.Duration `envconfig:"DB_CONN_TIMEOUT"`
		MaxConnTime time.Duration `envconfig:"DB_MAX_CONN_TIME"`
		MaxIdleTime time.Duration `envconfig:"DB_MAX_CONN_IDLE_TIME"`

		HealthCheckTime time.Duration `envconfig:"DB_HEALTH_CHECK_PERIOD"`
	}

	Redis struct {
		Address string `envconfig:"REDIS_ADDRESS"`

		DialTimeout  time.Duration `envconfig:"REDIS_DIAL_TIMEOUT"`
		ReadTimeout  time.Duration `envconfig:"REDIS_READ_TIMEOUT"`
		WriteTimeout time.Duration `envconfig:"REDIS_WRITE_TIMEOUT"`

		MaxRetries      int           `envconfig:"REDIS_MAX_RETRIES"`
		MinRetryBackoff time.Duration `envconfig:"REDIS_MIN_RETRY_BACKOFF"`
		MaxRetryBackoff time.Duration `envconfig:"REDIS_MAX_RETRY_BACKOFF"`

		PoolSize    int           `envconfig:"REDIS_POOL_SIZE"`
		MinIdleCons int           `envconfig:"REDIS_MIN_IDLE_CONS"`
		PoolTimeout time.Duration `envconfig:"REDIS_POOL_TIMEOUT"`
	}

	Otel struct {
		Endpoint string `envconfig:"OTEL_ENDPOINT"`
		IsTLS    bool   `envconfig:"OTEL_IS_TLS"`
	}

	Config struct {
		App    App
		Server Server
		DB     DB
		Redis  Redis
		Otel   Otel
	}
)

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("can't process the config: %w", err)
	}

	return &cfg, nil
}
