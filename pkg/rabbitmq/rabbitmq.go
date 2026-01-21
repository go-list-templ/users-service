package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/go-list-templ/grpc/config"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func New(cfg *config.RabbitMQ, logger *zap.Logger) (*RabbitMQ, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		logger.Warn("Failed to connect to RabbitMQ", zap.Error(err))
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Warn("Failed to open a RabbitMQ channel", zap.Error(err))
		return nil, err
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}
