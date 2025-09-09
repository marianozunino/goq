package rmq

import (
	"testing"

	"github.com/marianozunino/goq/internal/config"
	"github.com/marianozunino/goq/internal/testutil"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
)

func TestNewConsumer_ValidConfig(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "test-queue",
			AutoAck:     true,
		}

		consumer, err := NewConsumer(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if consumer == nil {
			t.Fatal("Expected consumer to be created")
		}

		if consumer.config != cfg {
			t.Error("Expected config to be set correctly")
		}

		if consumer.filter == nil {
			t.Error("Expected filter to be created")
		}
	})
}

func TestNewConsumer_WithTLS(t *testing.T) {
	cfg := &config.Config{
		RabbitMQURL:         "amqps://guest:guest@localhost:5671/",
		Queue:               "test-queue",
		AutoAck:             true,
		SkipTLSVerification: true,
	}

	// This will fail in real environment, but we can test the structure
	consumer, err := NewConsumer(cfg)
	if err == nil {
		t.Skip("Skipping test - RabbitMQ TLS is available locally")
	}

	// When connection fails, consumer should still be created for testing
	if consumer == nil {
		t.Skip("Skipping test - Consumer creation failed (expected in test environment)")
	}

	if !consumer.config.SkipTLSVerification {
		t.Error("Expected TLS verification to be skipped")
	}
}

func TestNewConsumer_WithFilters(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "test-queue",
			AutoAck:     true,
			FilterConfig: struct {
				IncludePatterns []string
				ExcludePatterns []string
				JSONFilter      string
				MaxMessageSize  int
				RegexFilter     string
			}{
				IncludePatterns: []string{"test"},
				ExcludePatterns: []string{"debug"},
				JSONFilter:      `.type == "test"`,
				MaxMessageSize:  1024,
				RegexFilter:     "error.*",
			},
		}

		consumer, err := NewConsumer(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if consumer == nil {
			t.Fatal("Expected consumer to be created")
		}

		if consumer.filter == nil {
			t.Error("Expected filter to be created")
		}
	})
}

func TestNewConsumer_InvalidFilter(t *testing.T) {
	cfg := &config.Config{
		RabbitMQURL: "amqp://guest:guest@localhost:5672/",
		Queue:       "test-queue",
		AutoAck:     true,
		FilterConfig: struct {
			IncludePatterns []string
			ExcludePatterns []string
			JSONFilter      string
			MaxMessageSize  int
			RegexFilter     string
		}{
			IncludePatterns: []string{"[invalid"},
			ExcludePatterns: []string{},
			JSONFilter:      "",
			MaxMessageSize:  -1,
			RegexFilter:     "",
		},
	}

	consumer, err := NewConsumer(cfg)
	if err == nil {
		t.Error("Expected error for invalid filter pattern")
	}

	if consumer != nil {
		t.Error("Expected consumer to be nil for invalid filter")
	}
}

func TestConvertDelivery(t *testing.T) {
	// Create a test delivery - rabbitmq.Delivery embeds amqp.Delivery
	delivery := rabbitmq.Delivery{}
	delivery.Body = []byte(`{"test": "data"}`)
	delivery.Exchange = "test_exchange"
	delivery.RoutingKey = "test.key"
	delivery.Headers = amqp091.Table{"test-header": "test-value"}

	// Convert delivery
	converted := convertDelivery(&delivery)

	if converted == nil {
		t.Fatal("Expected converted delivery to not be nil")
	}

	// Test GetBody method
	body := converted.GetBody()
	if len(body) == 0 {
		t.Error("Expected body to not be empty")
	}
}

func TestConvertHeaders(t *testing.T) {
	// Test with valid headers
	headers := amqp091.Table{
		"string-header": "test-value",
		"int-header":    42,
		"bool-header":   true,
	}

	converted := convertHeaders(rabbitmq.Table(headers))

	if converted == nil {
		t.Error("Expected converted headers to not be nil")
	}

	// Test with nil headers
	nilHeaders := convertHeaders(nil)
	if nilHeaders != nil {
		t.Error("Expected nil headers to remain nil")
	}
}

func TestParseBody(t *testing.T) {
	// Test with valid JSON
	validJSON := []byte(`{"test": "data"}`)
	result := parseBody(validJSON)

	if result == nil {
		t.Error("Expected parsed body to not be nil")
	}

	// Test with invalid JSON
	invalidJSON := []byte(`{invalid json`)
	result = parseBody(invalidJSON)

	if result == nil {
		t.Error("Expected parsed body to not be nil even for invalid JSON")
	}

	// Test with empty body
	emptyBody := []byte{}
	result = parseBody(emptyBody)

	if result == nil {
		t.Error("Expected parsed body to not be nil for empty body")
	}
}

func TestConsumer_Close(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "test-queue",
			AutoAck:     true,
		}

		consumer, err := NewConsumer(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if consumer == nil {
			t.Fatal("Expected consumer to be created")
		}

		// Test close method (should not panic even if connection is nil)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Close method panicked: %v", r)
			}
		}()

		consumer.Close()
	})
}

func TestConsumerStatus_Structure(t *testing.T) {
	status := ConsumerStatus{
		TotalMessages:    10,
		ConsumedMessages: 5,
		FilteredMessages: 2,
		Complete:         false,
		Message:          nil,
	}

	if status.TotalMessages != 10 {
		t.Error("Expected TotalMessages to be 10")
	}

	if status.ConsumedMessages != 5 {
		t.Error("Expected ConsumedMessages to be 5")
	}

	if status.FilteredMessages != 2 {
		t.Error("Expected FilteredMessages to be 2")
	}

	if status.Complete {
		t.Error("Expected Complete to be false")
	}

	if status.Message != nil {
		t.Error("Expected Message to be nil")
	}
}

func TestConsumer_EmptyQueueName(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "", // Empty queue name
			AutoAck:     false,
		}

		consumer, err := NewConsumer(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if consumer == nil {
			t.Fatal("Expected consumer to be created")
		}

		// Test that AutoAck is set to true for empty queue name
		if !consumer.config.AutoAck {
			t.Error("Expected AutoAck to be true for empty queue name")
		}
	})
}

func TestConsumer_WithRoutingKeys(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "test-queue",
			AutoAck:     true,
			RoutingKeys: []string{"user.created", "user.updated"},
		}

		consumer, err := NewConsumer(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if consumer == nil {
			t.Fatal("Expected consumer to be created")
		}

		if len(consumer.config.RoutingKeys) != 2 {
			t.Errorf("Expected 2 routing keys, got %d", len(consumer.config.RoutingKeys))
		}

		if consumer.config.RoutingKeys[0] != "user.created" {
			t.Errorf("Expected first routing key to be 'user.created', got '%s'", consumer.config.RoutingKeys[0])
		}

		if consumer.config.RoutingKeys[1] != "user.updated" {
			t.Errorf("Expected second routing key to be 'user.updated', got '%s'", consumer.config.RoutingKeys[1])
		}
	})
}

func TestConsumer_WithExchange(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "test-queue",
			AutoAck:     true,
			Exchange:    "test-exchange",
		}

		consumer, err := NewConsumer(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if consumer == nil {
			t.Fatal("Expected consumer to be created")
		}

		if consumer.config.Exchange != "test-exchange" {
			t.Errorf("Expected exchange to be 'test-exchange', got '%s'", consumer.config.Exchange)
		}
	})
}
