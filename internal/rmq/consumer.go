package rmq

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/marianozunino/goq/internal/config"
	"github.com/marianozunino/goq/internal/filter"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
)

type Consumer struct {
	conn     *rabbitmq.Conn
	consumer *rabbitmq.Consumer
	config   *config.Config
	filter   *filter.MessageFilter

	totalMessages    int
	consumedMessages int
}

type ConsumerStatus struct {
	TotalMessages    int
	ConsumedMessages int
	FilteredMessages int
	Complete         bool
	Message          *rabbitmq.Delivery
}

func NewConsumer(cfg *config.Config) (*Consumer, error) {
	msgFilter := filter.NewMessageFilter(cfg)
	if errs := msgFilter.GetCompilationErrors(); len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	// Create connection with TLS support
	var conn *rabbitmq.Conn
	var err error

	if cfg.SkipTLSVerification {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		conn, err = rabbitmq.NewConn(
			cfg.RabbitMQURL,
			rabbitmq.WithConnectionOptionsLogging,
			rabbitmq.WithConnectionOptionsConfig(rabbitmq.Config{
				TLSClientConfig: tlsConfig,
			}),
		)
	} else {
		conn, err = rabbitmq.NewConn(
			cfg.RabbitMQURL,
			rabbitmq.WithConnectionOptionsLogging,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	c := &Consumer{
		conn:   conn,
		config: cfg,
		filter: msgFilter,
	}

	// Handle queue setup
	if c.config.Queue == "" {
		c.config.AutoAck = true
	}

	return c, nil
}

// Consume starts consuming messages and returns a channel for status updates
func (c *Consumer) Consume() (<-chan ConsumerStatus, error) {
	statusCh := make(chan ConsumerStatus)

	// Prepare consumer options
	var consumerOptions []func(*rabbitmq.ConsumerOptions)

	// Handle queue creation - if no queue name is provided, create a temporary queue
	queueName := c.config.Queue
	if queueName == "" {
		consumerOptions = append(consumerOptions,
			rabbitmq.WithConsumerOptionsQueueAutoDelete,
			rabbitmq.WithConsumerOptionsQueueExclusive,
		)
		queueName = ""
	} else {
		// For named queues, make them durable to match existing queue configurations
		consumerOptions = append(consumerOptions,
			rabbitmq.WithConsumerOptionsQueueDurable,
		)
	}

	// Add routing keys if specified
	for _, routingKey := range c.config.RoutingKeys {
		consumerOptions = append(consumerOptions,
			rabbitmq.WithConsumerOptionsRoutingKey(routingKey))
	}

	// Add exchange configuration
	if c.config.Exchange != "" {
		consumerOptions = append(consumerOptions,
			rabbitmq.WithConsumerOptionsExchangeName(c.config.Exchange),
		)
	}

	// Create consumer
	consumer, err := rabbitmq.NewConsumer(
		c.conn,
		queueName,
		consumerOptions...,
	)
	if err != nil {
		close(statusCh)
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	c.consumer = consumer

	// Get queue message count after consumer is created (for no-ack mode with StopAfterConsume)
	if !c.config.AutoAck && c.config.StopAfterConsume && queueName != "" {
		// Get queue info to determine message count using amqp091-go directly
		// Use the same TLS configuration as the main connection
		var conn *amqp091.Connection
		var err error

		if c.config.SkipTLSVerification {
			// Configure TLS for secure connections with skip verification
			tlsConfig := &tls.Config{InsecureSkipVerify: true}
			conn, err = amqp091.DialTLS(c.config.RabbitMQURL, tlsConfig)
		} else {
			conn, err = amqp091.Dial(c.config.RabbitMQURL)
		}

		if err != nil {
			close(statusCh)
			return nil, fmt.Errorf("failed to connect for queue info: %v", err)
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			close(statusCh)
			return nil, fmt.Errorf("failed to get channel for queue info: %v", err)
		}
		defer ch.Close()

		queue, err := ch.QueueInspect(queueName)
		if err != nil {
			// Queue might not exist yet, that's okay - we'll set totalMessages to 0
			c.totalMessages = 0
			fmt.Printf("📊 Queue %s not found or empty\n", queueName)
		} else {
			c.totalMessages = queue.Messages
			fmt.Printf("📊 Queue has %d messages\n", c.totalMessages)
		}
	}

	if queueName == "" {
		fmt.Println("✅ Temporary queue created with random name (managed by go-rabbitmq)")
	} else {
		fmt.Printf("✅ Connected to queue: %s\n", queueName)
	}

	if len(c.config.RoutingKeys) > 0 {
		fmt.Printf("✅ Bound to routing keys: %v\n", c.config.RoutingKeys)
	}
	if c.config.Exchange != "" {
		fmt.Printf("✅ Connected to exchange: %s\n", c.config.Exchange)
	}

	go func() {
		defer close(statusCh)
		filteredCount := 0
		messageCount := 0

		err := consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
			c.consumedMessages++
			messageCount++
			var filteredMsg *rabbitmq.Delivery

			if c.filter.Filter(convertDelivery(&d)) {
				filteredMsg = &d
			} else {
				filteredCount++
			}

			statusCh <- ConsumerStatus{
				TotalMessages:    c.totalMessages,
				ConsumedMessages: c.consumedMessages,
				FilteredMessages: filteredCount,
				Complete:         false,
				Message:          filteredMsg,
			}

			// For no-ack mode with StopAfterConsume, stop when we've processed enough messages
			if !c.config.AutoAck && c.config.StopAfterConsume && c.totalMessages > 0 && messageCount >= c.totalMessages {
				statusCh <- ConsumerStatus{
					TotalMessages:    c.totalMessages,
					ConsumedMessages: c.consumedMessages,
					FilteredMessages: filteredCount,
					Complete:         true,
					Message:          nil,
				}
				return rabbitmq.NackRequeue
			}

			if c.config.AutoAck {
				return rabbitmq.Ack
			}
			return rabbitmq.NackRequeue
		})

		if err != nil {
			fmt.Printf("Consumer error: %v\n", err)
		}
	}()

	return statusCh, nil
}

// Close closes the consumer and connection
func (c *Consumer) Close() error {
	if c.consumer != nil {
		c.consumer.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}

// convertDelivery converts rabbitmq.Delivery to amqp.Delivery for compatibility with existing filter
func convertDelivery(d *rabbitmq.Delivery) *amqpDelivery {
	message := struct {
		Headers    map[string]interface{} `json:"headers"`
		Exchange   string                 `json:"exchange"`
		RoutingKey string                 `json:"routingKey"`
		Body       interface{}            `json:"body"`
	}{
		Headers:    convertHeaders(rabbitmq.Table(d.Headers)),
		Exchange:   d.Exchange,
		RoutingKey: d.RoutingKey,
		Body:       parseBody(d.Body),
	}

	jsonBytes, _ := json.Marshal(message)

	return &amqpDelivery{
		Body: jsonBytes,
	}
}

// convertHeaders converts rabbitmq.Table to map[string]interface{}
func convertHeaders(headers rabbitmq.Table) map[string]interface{} {
	if headers == nil {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range headers {
		result[k] = v
	}
	return result
}

// parseBody attempts to parse the body as JSON, falls back to string
func parseBody(body []byte) interface{} {
	var jsonData interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return string(body)
	}
	return jsonData
}

// amqpDelivery is a minimal adapter to make rabbitmq.Delivery compatible with the existing filter
type amqpDelivery struct {
	Body []byte
}

// GetBody implements the MessageDelivery interface
func (d *amqpDelivery) GetBody() []byte {
	return d.Body
}
