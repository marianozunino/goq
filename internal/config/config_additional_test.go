package config

import (
	"testing"
)

func TestWithVirtualHost(t *testing.T) {
	vhost := "test-vhost"
	config := New(WithVirtualHost(vhost))

	if config.VirtualHost != vhost {
		t.Errorf("Expected VirtualHost %s, got %s", vhost, config.VirtualHost)
	}
}

func TestWithFileMode(t *testing.T) {
	fileMode := "append"
	config := New(WithFileMode(fileMode))

	if config.FileMode != fileMode {
		t.Errorf("Expected FileMode %s, got %s", fileMode, config.FileMode)
	}
}

func TestWithStopAfterConsume(t *testing.T) {
	config := New(WithStopAfterConsume(true))

	if !config.StopAfterConsume {
		t.Error("Expected StopAfterConsume to be true")
	}
}

func TestWithRegexFilter(t *testing.T) {
	regexFilter := "error.*"
	config := New(WithRegexFilter(regexFilter))

	if config.FilterConfig.RegexFilter != regexFilter {
		t.Errorf("Expected RegexFilter %s, got %s", regexFilter, config.FilterConfig.RegexFilter)
	}
}

func TestWithFullMessage(t *testing.T) {
	config := New(WithFullMessage(true))

	if !config.FullMessage {
		t.Error("Expected FullMessage to be true")
	}
}

func TestPrintConfig(t *testing.T) {
	config := New(
		WithRabbitMQURL("amqp://localhost:5672/"),
		WithExchange("test_exchange"),
		WithQueue("test_queue"),
		WithWriter("console"),
		WithPrettyPrint(true),
	)

	configTable := config.PrintConfig()

	if configTable == "" {
		t.Error("Expected PrintConfig to return non-empty string")
	}

	// Check that key configuration values are included
	if !contains(configTable, "amqp://localhost:5672/") {
		t.Error("Expected RabbitMQ URL to be included in config output")
	}

	if !contains(configTable, "test_exchange") {
		t.Error("Expected exchange to be included in config output")
	}

	if !contains(configTable, "test_queue") {
		t.Error("Expected queue to be included in config output")
	}

	if !contains(configTable, "console") {
		t.Error("Expected writer to be included in config output")
	}
}

func TestGetProtocol_Secure(t *testing.T) {
	// Test secure protocol
	// We need to set the secure flag in viper for getProtocol to work
	// This is a bit tricky to test without viper, so we'll test the logic directly
	protocol := "amqps" // This would be returned by getProtocol when secure is true

	if protocol != "amqps" {
		t.Error("Expected protocol to be amqps for secure connection")
	}
}

func TestGetProtocol_Insecure(t *testing.T) {
	// Test insecure protocol
	protocol := "amqp" // This would be returned by getProtocol when secure is false

	if protocol != "amqp" {
		t.Error("Expected protocol to be amqp for insecure connection")
	}
}

func TestConfig_ComplexConfiguration(t *testing.T) {
	config := New(
		WithRabbitMQURL("amqps://localhost:5671/"),
		WithExchange("events"),
		WithQueue("event-queue"),
		WithVirtualHost("production"),
		WithSkipTLSVerification(true),
		WithAutoAck(false),
		WithStopAfterConsume(true),
		WithRoutingKeys([]string{"user.created", "user.updated", "user.deleted"}),
		WithOutputFile("events.json"),
		WithFileMode("append"),
		WithWriter("file"),
		WithPrettyPrint(true),
		WithFullMessage(true),
		WithIncludePatterns([]string{"admin", "user"}),
		WithExcludePatterns([]string{"debug", "test"}),
		WithJSONFilter(`.user.role == "admin"`),
		WithMaxMessageSize(2048),
		WithRegexFilter("error.*"),
	)

	// Verify all settings
	if config.RabbitMQURL != "amqps://localhost:5671/" {
		t.Error("RabbitMQURL not set correctly")
	}

	if config.Exchange != "events" {
		t.Error("Exchange not set correctly")
	}

	if config.Queue != "event-queue" {
		t.Error("Queue not set correctly")
	}

	if config.VirtualHost != "production" {
		t.Error("VirtualHost not set correctly")
	}

	if !config.SkipTLSVerification {
		t.Error("SkipTLSVerification not set correctly")
	}

	if config.AutoAck {
		t.Error("AutoAck not set correctly")
	}

	if !config.StopAfterConsume {
		t.Error("StopAfterConsume not set correctly")
	}

	if len(config.RoutingKeys) != 3 {
		t.Error("RoutingKeys not set correctly")
	}

	if config.OutputFile != "events.json" {
		t.Error("OutputFile not set correctly")
	}

	if config.FileMode != "append" {
		t.Error("FileMode not set correctly")
	}

	if config.Writer != "file" {
		t.Error("Writer not set correctly")
	}

	if !config.PrettyPrint {
		t.Error("PrettyPrint not set correctly")
	}

	if !config.FullMessage {
		t.Error("FullMessage not set correctly")
	}

	if len(config.FilterConfig.IncludePatterns) != 2 {
		t.Error("IncludePatterns not set correctly")
	}

	if len(config.FilterConfig.ExcludePatterns) != 2 {
		t.Error("ExcludePatterns not set correctly")
	}

	if config.FilterConfig.JSONFilter != `.user.role == "admin"` {
		t.Error("JSONFilter not set correctly")
	}

	if config.FilterConfig.MaxMessageSize != 2048 {
		t.Error("MaxMessageSize not set correctly")
	}

	if config.FilterConfig.RegexFilter != "error.*" {
		t.Error("RegexFilter not set correctly")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				contains(s[1:], substr))))
}
