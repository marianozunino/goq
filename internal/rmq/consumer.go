package rmq

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/marianozunino/goq/internal/config"
	"github.com/streadway/amqp"
)

type Message = amqp.Delivery

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

	filteredMsgs := make(chan amqp.Delivery)

	go func() {
		defer close(filteredMsgs)

		for msg := range msgs {
			if c.shouldProcessMessage(&msg) {
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

// DeclareTemporaryQueue creates a temporary queue that will be deleted when the consumer disconnects
func (c *Consumer) DeclareTemporaryQueue() (amqp.Queue, error) {
	return c.ch.QueueDeclare(
		"",    // Generate a random name for the queue
		false, // Durable
		true,  // Delete when unused
		true,  // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
}

// BindQueue binds a queue to an exchange using a routing key
func (c *Consumer) BindQueue(queueName, routingKey string) error {
	return c.ch.QueueBind(
		queueName,
		routingKey,
		c.config.Exchange,
		false,
		nil,
	)
}

// ConsumeMessagesFromQueue consumes messages from the specified queue
func (c *Consumer) ConsumeMessagesFromQueue(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := c.ch.Consume(
		queueName,
		"",    // Consumer tag
		false, // Auto-ack
		false, // Exclusive
		false, // No wait
		false, // No local
		nil,   // Arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer for queue %s: %v", queueName, err)
	}

	filteredMsgs := make(chan amqp.Delivery)

	go func() {
		defer close(filteredMsgs)

		for msg := range msgs {
			if c.shouldProcessMessage(&msg) {
				filteredMsgs <- msg
			}
		}
	}()

	return filteredMsgs, nil
}

func (c *Consumer) shouldProcessMessage(msg *Message) bool {
	config := c.config.FilterConfig

	// Size filtering
	if config.MaxMessageSize > 0 && len(msg.Body) > config.MaxMessageSize {
		return false
	}

	// Regex filtering
	if config.CompileRegex != nil && !config.CompileRegex.Match(msg.Body) {
		return false
	}

	body := string(msg.Body)

	// Include patterns
	if len(config.IncludePatterns) > 0 {
		if !sliceContainsAny(config.IncludePatterns, body) {
			return false
		}
	}

	// Exclude patterns
	if sliceContainsAny(config.ExcludePatterns, body) {
		return false
	}

	// JSON filter
	if config.JSONFilter != nil {
		return matchJSONFilter(msg.Body, config.JSONFilter)
	}

	return true
}

func sliceContainsAny(patterns []string, body string) bool {
	for _, pattern := range patterns {
		if strings.Contains(body, pattern) {
			return true
		}
	}
	return false
}

func matchJSONFilter(body []byte, query *gojq.Query) bool {
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return false
	}

	iter := query.Run(data)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, isErr := v.(error); isErr {
			fmt.Printf("Error applying jq filter: %v\n", err)
			return false
		}

		// Check if result is truthy
		switch val := v.(type) {
		case bool:
			return val
		case nil:
			continue
		default:
			result, _ := json.MarshalIndent(v, "", "  ")
			return strings.TrimSpace(string(result)) != ""
		}
	}
	return false
}
