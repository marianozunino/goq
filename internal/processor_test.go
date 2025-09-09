package app

import (
	"testing"

	"github.com/marianozunino/goq/internal/config"
	"github.com/marianozunino/goq/internal/testutil"
)

func TestNewMessageProcessor_ValidConfig(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "test-queue",
			AutoAck:     true,
			Writer:      "console",
		}

		processor, err := NewMessageProcessor(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if processor == nil {
			t.Fatal("Expected processor to be created")
		}

		if processor.config != cfg {
			t.Error("Expected config to be set correctly")
		}
	})
}

func TestNewMessageProcessor_WithFileExporter(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "test-queue",
			AutoAck:     true,
			Writer:      "file",
			OutputFile:  "test.json",
			FileMode:    "overwrite",
		}

		processor, err := NewMessageProcessor(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if processor == nil {
			t.Fatal("Expected processor to be created")
		}

		if processor.config.Writer != "file" {
			t.Error("Expected writer to be set to file")
		}

		if processor.config.OutputFile != "test.json" {
			t.Error("Expected output file to be set correctly")
		}
	})
}

func TestNewMessageProcessor_WithFilters(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "test-queue",
			AutoAck:     true,
			Writer:      "console",
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

		processor, err := NewMessageProcessor(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if processor == nil {
			t.Fatal("Expected processor to be created")
		}

		if len(processor.config.FilterConfig.IncludePatterns) != 1 {
			t.Error("Expected include patterns to be set")
		}

		if len(processor.config.FilterConfig.ExcludePatterns) != 1 {
			t.Error("Expected exclude patterns to be set")
		}

		if processor.config.FilterConfig.JSONFilter != `.type == "test"` {
			t.Error("Expected JSON filter to be set")
		}

		if processor.config.FilterConfig.MaxMessageSize != 1024 {
			t.Error("Expected max message size to be set")
		}

		if processor.config.FilterConfig.RegexFilter != "error.*" {
			t.Error("Expected regex filter to be set")
		}
	})
}

func TestNewMessageProcessor_InvalidExporter(t *testing.T) {
	cfg := &config.Config{
		RabbitMQURL: "amqp://guest:guest@localhost:5672/",
		Queue:       "test-queue",
		AutoAck:     true,
		Writer:      "invalid-writer",
	}

	processor, err := NewMessageProcessor(cfg)
	if err == nil {
		t.Error("Expected error for invalid writer type")
	}

	if processor != nil {
		t.Error("Expected processor to be nil for invalid writer")
	}
}

func TestNewMessageProcessor_FileWriterWithoutOutput(t *testing.T) {
	cfg := &config.Config{
		RabbitMQURL: "amqp://guest:guest@localhost:5672/",
		Queue:       "test-queue",
		AutoAck:     true,
		Writer:      "file",
		OutputFile:  "", // Empty output file
		FileMode:    "overwrite",
	}

	processor, err := NewMessageProcessor(cfg)
	if err == nil {
		t.Error("Expected error for file writer without output file")
	}

	if processor != nil {
		t.Error("Expected processor to be nil for file writer without output")
	}
}

func TestMessageProcessor_Structure(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "test-queue",
			AutoAck:     true,
			Writer:      "console",
		}

		processor, err := NewMessageProcessor(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if processor == nil {
			t.Fatal("Expected processor to be created")
		}

		// Test that all components are set
		if processor.config == nil {
			t.Error("Expected config to be set")
		}

		if processor.consumer == nil {
			t.Error("Expected consumer to be set")
		}

		if processor.exporter == nil {
			t.Error("Expected exporter to be set")
		}
	})
}

func TestMessageProcessor_WithPrettyPrint(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "test-queue",
			AutoAck:     true,
			Writer:      "console",
			PrettyPrint: true,
		}

		processor, err := NewMessageProcessor(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if processor == nil {
			t.Fatal("Expected processor to be created")
		}

		if !processor.config.PrettyPrint {
			t.Error("Expected pretty print to be enabled")
		}
	})
}

func TestMessageProcessor_WithFullMessage(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "test-queue",
			AutoAck:     true,
			Writer:      "console",
			FullMessage: true,
		}

		processor, err := NewMessageProcessor(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if processor == nil {
			t.Fatal("Expected processor to be created")
		}

		if !processor.config.FullMessage {
			t.Error("Expected full message to be enabled")
		}
	})
}

func TestMessageProcessor_WithStopAfterConsume(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL:      rmq.GetConnectionURL(),
			Queue:            "test-queue",
			AutoAck:          true,
			Writer:           "console",
			StopAfterConsume: true,
		}

		processor, err := NewMessageProcessor(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if processor == nil {
			t.Fatal("Expected processor to be created")
		}

		if !processor.config.StopAfterConsume {
			t.Error("Expected stop after consume to be enabled")
		}
	})
}

func TestMessageProcessor_WithRoutingKeys(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "test-queue",
			AutoAck:     true,
			Writer:      "console",
			RoutingKeys: []string{"user.created", "user.updated"},
		}

		processor, err := NewMessageProcessor(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if processor == nil {
			t.Fatal("Expected processor to be created")
		}

		if len(processor.config.RoutingKeys) != 2 {
			t.Errorf("Expected 2 routing keys, got %d", len(processor.config.RoutingKeys))
		}
	})
}

func TestMessageProcessor_WithExchange(t *testing.T) {
	testutil.WithRabbitMQTestContainer(t, func(rmq *testutil.RabbitMQTestContainer) {
		cfg := &config.Config{
			RabbitMQURL: rmq.GetConnectionURL(),
			Queue:       "test-queue",
			AutoAck:     true,
			Writer:      "console",
			Exchange:    "test-exchange",
		}

		processor, err := NewMessageProcessor(cfg)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if processor == nil {
			t.Fatal("Expected processor to be created")
		}

		if processor.config.Exchange != "test-exchange" {
			t.Errorf("Expected exchange to be 'test-exchange', got '%s'", processor.config.Exchange)
		}
	})
}
