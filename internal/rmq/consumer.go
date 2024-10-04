package rmq

import (
	"crypto/tls"
	"fmt"

	"github.com/marianozunino/goq/internal/config"
	"github.com/streadway/amqp"
)

type Consumer struct {
	config *config.Config
	conn   *amqp.Connection
	ch     *amqp.Channel
}

func NewConsumer(cfg *config.Config) (*Consumer, error) {
	var conn *amqp.Connection
	var err error

	if cfg.SkipTLSVerification && cfg.UseAMQPS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		conn, err = amqp.DialTLS(cfg.RabbitMQURL, tlsConfig)
	} else {
		conn, err = amqp.Dial(cfg.RabbitMQURL)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	return &Consumer{
		config: cfg,
		conn:   conn,
		ch:     ch,
	}, nil
}

func (c *Consumer) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Consumer) ConsumeMessages() (<-chan amqp.Delivery, error) {
	err := c.ch.QueueBind(
		c.config.Queue,
		"",
		c.config.Exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind a queue: %v", err)
	}

	msgs, err := c.ch.Consume(
		c.config.Queue,
		"",
		c.config.AutoAck,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %v", err)
	}

	return msgs, nil
}

func (c *Consumer) GetQueueInfo() (int, error) {
	queue, err := c.ch.QueueInspect(c.config.Queue)
	if err != nil {
		return 0, fmt.Errorf("failed to inspect queue: %v", err)
	}
	return queue.Messages, nil
}