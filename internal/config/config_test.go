package config

import (
	"testing"
)

func TestNewConfig_DefaultValues(t *testing.T) {
	config := New()

	// These should be false by default
	if config.AutoAck {
		t.Error("Expected AutoAck to be false by default")
	}

	if config.SkipTLSVerification {
		t.Error("Expected SkipTLSVerification to be false by default")
	}

	if config.PrettyPrint {
		t.Error("Expected PrettyPrint to be false by default")
	}

	if config.FullMessage {
		t.Error("Expected FullMessage to be false by default")
	}
}

func TestWithRabbitMQURL(t *testing.T) {
	url := "amqp://localhost:5672/"
	config := New(WithRabbitMQURL(url))

	if config.RabbitMQURL != url {
		t.Errorf("Expected RabbitMQURL %s, got %s", url, config.RabbitMQURL)
	}
}

func TestWithExchange(t *testing.T) {
	exchange := "test_exchange"
	config := New(WithExchange(exchange))

	if config.Exchange != exchange {
		t.Errorf("Expected Exchange %s, got %s", exchange, config.Exchange)
	}
}

func TestWithQueue(t *testing.T) {
	queue := "test_queue"
	config := New(WithQueue(queue))

	if config.Queue != queue {
		t.Errorf("Expected Queue %s, got %s", queue, config.Queue)
	}
}

func TestWithRoutingKeys(t *testing.T) {
	routingKeys := []string{"user.created", "user.updated"}
	config := New(WithRoutingKeys(routingKeys))

	if len(config.RoutingKeys) != len(routingKeys) {
		t.Errorf("Expected %d routing keys, got %d", len(routingKeys), len(config.RoutingKeys))
	}

	for i, key := range routingKeys {
		if config.RoutingKeys[i] != key {
			t.Errorf("Expected routing key %s at index %d, got %s", key, i, config.RoutingKeys[i])
		}
	}
}

func TestWithAutoAck(t *testing.T) {
	config := New(WithAutoAck(true))

	if !config.AutoAck {
		t.Error("Expected AutoAck to be true")
	}
}

func TestWithSkipTLSVerification(t *testing.T) {
	config := New(WithSkipTLSVerification(true))

	if !config.SkipTLSVerification {
		t.Error("Expected SkipTLSVerification to be true")
	}
}

func TestWithOutputFile(t *testing.T) {
	outputFile := "test_output.json"
	config := New(WithOutputFile(outputFile))

	if config.OutputFile != outputFile {
		t.Errorf("Expected OutputFile %s, got %s", outputFile, config.OutputFile)
	}
}

func TestWithWriter(t *testing.T) {
	writer := "console"
	config := New(WithWriter(writer))

	if config.Writer != ExporterKind(writer) {
		t.Errorf("Expected Writer %s, got %s", writer, config.Writer)
	}
}

func TestWithPrettyPrint(t *testing.T) {
	config := New(WithPrettyPrint(true))

	if !config.PrettyPrint {
		t.Error("Expected PrettyPrint to be true")
	}
}

func TestWithIncludePatterns(t *testing.T) {
	patterns := []string{"admin", "error"}
	config := New(WithIncludePatterns(patterns))

	if len(config.FilterConfig.IncludePatterns) != len(patterns) {
		t.Errorf("Expected %d include patterns, got %d", len(patterns), len(config.FilterConfig.IncludePatterns))
	}

	for i, pattern := range patterns {
		if config.FilterConfig.IncludePatterns[i] != pattern {
			t.Errorf("Expected include pattern %s at index %d, got %s", pattern, i, config.FilterConfig.IncludePatterns[i])
		}
	}
}

func TestWithExcludePatterns(t *testing.T) {
	patterns := []string{"debug", "test"}
	config := New(WithExcludePatterns(patterns))

	if len(config.FilterConfig.ExcludePatterns) != len(patterns) {
		t.Errorf("Expected %d exclude patterns, got %d", len(patterns), len(config.FilterConfig.ExcludePatterns))
	}

	for i, pattern := range patterns {
		if config.FilterConfig.ExcludePatterns[i] != pattern {
			t.Errorf("Expected exclude pattern %s at index %d, got %s", pattern, i, config.FilterConfig.ExcludePatterns[i])
		}
	}
}

func TestWithJSONFilter(t *testing.T) {
	jsonFilter := `.user.role == "admin"`
	config := New(WithJSONFilter(jsonFilter))

	if config.FilterConfig.JSONFilter != jsonFilter {
		t.Errorf("Expected JSONFilter %s, got %s", jsonFilter, config.FilterConfig.JSONFilter)
	}
}

func TestWithMaxMessageSize(t *testing.T) {
	maxSize := 1024
	config := New(WithMaxMessageSize(maxSize))

	if config.FilterConfig.MaxMessageSize != maxSize {
		t.Errorf("Expected MaxMessageSize %d, got %d", maxSize, config.FilterConfig.MaxMessageSize)
	}
}

func TestMultipleOptions(t *testing.T) {
	config := New(
		WithRabbitMQURL("amqp://localhost:5672/"),
		WithExchange("test_exchange"),
		WithQueue("test_queue"),
		WithAutoAck(true),
		WithPrettyPrint(true),
	)

	if config.RabbitMQURL != "amqp://localhost:5672/" {
		t.Error("RabbitMQURL not set correctly")
	}

	if config.Exchange != "test_exchange" {
		t.Error("Exchange not set correctly")
	}

	if config.Queue != "test_queue" {
		t.Error("Queue not set correctly")
	}

	if !config.AutoAck {
		t.Error("AutoAck not set correctly")
	}

	if !config.PrettyPrint {
		t.Error("PrettyPrint not set correctly")
	}
}
