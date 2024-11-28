package rmq

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/marianozunino/goq/internal/config"
	"github.com/marianozunino/goq/internal/filter"
	"github.com/streadway/amqp"
)

type Consumer struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	config *config.Config
	filter *filter.MessageFilter

	totalMessages    int
	consumedMessages int
}

type ConsumerStatus struct {
	TotalMessages    int
	ConsumedMessages int
	FilteredMessages int
	Complete         bool
	Message          *amqp.Delivery
}

func NewConsumer(cfg *config.Config) (*Consumer, error) {
	msgFilter := filter.NewMessageFilter(cfg)
	if errs := msgFilter.GetCompilationErrors(); len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

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

	c := &Consumer{
		conn:   conn,
		ch:     ch,
		config: cfg,
		filter: msgFilter,
	}

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

	if err = c.ch.QueueBind(c.config.Queue, "", c.config.Exchange, false, nil); err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	queue, err := c.ch.QueueInspect(c.config.Queue)
	c.totalMessages = int(queue.Messages)

	if err != nil {
		return nil, fmt.Errorf("failed to inspect queue: %v", err)
	}

	return c, nil
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

func (c *Consumer) Consume() (<-chan ConsumerStatus, error) {
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
		return nil, fmt.Errorf("failed to register consumer: %v", err)
	}
	statusCh := make(chan ConsumerStatus)
	go func(c *Consumer) {
		defer close(statusCh)
		filteredCount := 0
		for msg := range msgs {
			c.consumedMessages++
			var filteredMsg *amqp.Delivery
			if c.filter.Filter(&msg) {
				filteredMsg = &msg
			} else {
				filteredCount++
			}

			// Send status update
			statusCh <- ConsumerStatus{
				TotalMessages:    c.totalMessages,
				ConsumedMessages: c.consumedMessages,
				FilteredMessages: filteredCount,
				Complete:         c.consumedMessages == c.totalMessages,
				Message:          filteredMsg,
			}

			if c.consumedMessages == c.totalMessages {
				break
			}
		}
	}(c)
	return statusCh, nil
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
