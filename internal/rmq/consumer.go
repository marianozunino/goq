package rmq

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/marianozunino/goq/internal/config"
	"github.com/streadway/amqp"
)

type Consumer struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	config *config.Config
}

func NewConsumer(cfg *config.Config) (*Consumer, error) {
	var conn *amqp.Connection
	var err error

	if cfg.SkipTLSVerification {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
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
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	return &Consumer{
		conn:   conn,
		ch:     ch,
		config: cfg,
	}, nil
}

// DeclareTemporaryQueue creates a temporary queue that will be deleted when the consumer disconnects
func (c *Consumer) declareTemporaryQueue() (amqp.Queue, error) {
	fmt.Println("Declaring a temporary queue")
	return c.ch.QueueDeclare(
		"",    // Generate a random name for the queue
		false, // Durable
		true,  // Delete when unused
		true,  // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	if c.config.Queue == "" {
		queue, err := c.declareTemporaryQueue()
		if err != nil {
			return nil, fmt.Errorf("failed to declare temporary queue: %v", err)
		}
		c.config.Queue = queue.Name
		c.config.AutoAck = true

		err = c.bindQueueToRoutingKeys(queue.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to bind queue to routing keys: %v", err)
		}
	}

	err := c.ch.QueueBind(c.config.Queue, "", c.config.Exchange, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	msgs, err := c.ch.Consume(
		c.config.Queue,   // queue
		"",               // consumer
		c.config.AutoAck, // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register consumer: %v", err)
	}

	color.Cyan("Monitoring messages on temporary queue: %s", c.config.Queue)

	filteredMsgs := make(chan amqp.Delivery)

	go func() {
		defer close(filteredMsgs)
		for msg := range msgs {
			if c.filterMessage(&msg) {
				filteredMsgs <- msg
			}
		}
	}()

	return filteredMsgs, nil
}

func (c *Consumer) GetQueueInfo() (int, error) {
	queue, err := c.ch.QueueInspect(c.config.Queue)
	if err != nil {
		return 0, fmt.Errorf("failed to inspect queue: %v", err)
	}
	return queue.Messages, nil
}

// bindQueueToRoutingKeys binds the queue to specified routing keys
func (c *Consumer) bindQueueToRoutingKeys(queueName string) error {
	for _, routingKey := range c.config.RoutingKeys {
		routingKeyTrimmed := strings.TrimSpace(routingKey)

		if err := c.ch.QueueBind(
			queueName,
			routingKey,
			c.config.Exchange,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to bind queue %q to routing key %q: %v", queueName, routingKeyTrimmed, err)
		}

		color.Yellow("Bound routing key %q", routingKey)
		color.Green("Binding queue %q to routing keys", queueName)
	}
	return nil
}
