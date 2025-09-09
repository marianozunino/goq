package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/wait"
)

// RabbitMQTestContainer represents a test RabbitMQ container
type RabbitMQTestContainer struct {
	container testcontainers.Container
	host      string
	port      string
	username  string
	password  string
	amqpURL   string
}

// NewRabbitMQTestContainer creates a new RabbitMQ test container
func NewRabbitMQTestContainer(ctx context.Context, t *testing.T) (*RabbitMQTestContainer, error) {
	// Create RabbitMQ container
	rabbitmqContainer, err := rabbitmq.RunContainer(ctx,
		testcontainers.WithImage("rabbitmq:3.13-management"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Server startup complete").
				WithOccurrence(1).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start RabbitMQ container: %w", err)
	}

	// Get connection details
	host, err := rabbitmqContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host: %w", err)
	}

	port, err := rabbitmqContainer.MappedPort(ctx, "5672")
	if err != nil {
		return nil, fmt.Errorf("failed to get port: %w", err)
	}

	// Get AMQP URL directly from the container
	amqpURL, err := rabbitmqContainer.AmqpURL(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get AMQP URL: %w", err)
	}

	return &RabbitMQTestContainer{
		container: rabbitmqContainer,
		host:      host,
		port:      port.Port(),
		username:  rabbitmqContainer.AdminUsername,
		password:  rabbitmqContainer.AdminPassword,
		amqpURL:   amqpURL,
	}, nil
}

// GetConnectionURL returns the AMQP connection URL
func (r *RabbitMQTestContainer) GetConnectionURL() string {
	return r.amqpURL
}

// GetManagementURL returns the management UI URL
func (r *RabbitMQTestContainer) GetManagementURL() string {
	return fmt.Sprintf("http://%s:15672", r.host)
}

// Close terminates the container
func (r *RabbitMQTestContainer) Close(ctx context.Context) error {
	return r.container.Terminate(ctx)
}

// Cleanup is a convenience method for tests that automatically handles cleanup
func (r *RabbitMQTestContainer) Cleanup(t *testing.T) {
	ctx := context.Background()
	if err := r.Close(ctx); err != nil {
		t.Logf("Failed to cleanup RabbitMQ container: %v", err)
	}
}

// WithRabbitMQTestContainer is a helper function that sets up a RabbitMQ container for testing
func WithRabbitMQTestContainer(t *testing.T, testFunc func(*RabbitMQTestContainer)) {
	ctx := context.Background()

	container, err := NewRabbitMQTestContainer(ctx, t)
	if err != nil {
		t.Fatalf("Failed to create RabbitMQ test container: %v", err)
	}
	defer container.Cleanup(t)

	// Wait a bit for RabbitMQ to be fully ready
	time.Sleep(2 * time.Second)

	testFunc(container)
}
