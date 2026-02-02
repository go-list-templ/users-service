package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	DB struct {
		URL         string        `envconfig:"DB_URL"`
		Driver      string        `envconfig:"DB_DRIVER"`
		MaxConn     int32         `envconfig:"DB_MAX_CONN"`
		MaxIdle     int           `envconfig:"DB_MAX_IDLE"`
		MaxConnTime time.Duration `envconfig:"DB_MAX_CONN_TIME"`
		MaxIdleTime time.Duration `envconfig:"DB_MAX_IDLE_TIME"`
	}

	Redis struct {
		Address string `envconfig:"REDIS_ADDRESS"`
	}

	App struct {
		Name    string `envconfig:"APP_NAME"`
		Version string `envconfig:"APP_VERSION"`
	}

	Server struct {
		GRPCPort        string        `envconfig:"GRPC_PORT"`
		HealthPort      string        `envconfig:"HEALTH_PORT"`
		HTTPTimeout     time.Duration `envconfig:"HTTP_TIMEOUT"`
		IdleTimeout     time.Duration `envconfig:"IDLE_TIMEOUT"`
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT"`
	}

	Config struct {
		App    App
		Server Server
		DB     DB
		Redis  Redis
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
