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
		DiagnosticPort  string        `envconfig:"DIAGNOSTIC_PORT"`
		HTTPTimeout     time.Duration `envconfig:"HTTP_TIMEOUT"`
		IdleTimeout     time.Duration `envconfig:"IDLE_TIMEOUT"`
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT"`
	}

	DB struct {
		ReplicaURL string `envconfig:"DB_REPLICA_URL"`
		MasterURL  string `envconfig:"DB_MASTER_URL"`

		MaxConn     int32         `envconfig:"DB_MAX_CONN"`
		MaxConnTime time.Duration `envconfig:"DB_MAX_CONN_TIME"`
		MaxIdleTime time.Duration `envconfig:"DB_MAX_CONN_IDLE_TIME"`
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
