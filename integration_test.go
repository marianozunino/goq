package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marianozunino/goq/internal/config"
	"github.com/marianozunino/goq/internal/exporter"
	"github.com/marianozunino/goq/internal/filter"
)

func TestIntegration_FilterAndExport(t *testing.T) {
	// Create a test configuration
	cfg := &config.Config{
		FilterConfig: struct {
			IncludePatterns []string
			ExcludePatterns []string
			JSONFilter      string
			MaxMessageSize  int
			RegexFilter     string
		}{
			IncludePatterns: []string{"admin"},
			ExcludePatterns: []string{"debug"},
			JSONFilter:      "",
			MaxMessageSize:  -1,
			RegexFilter:     "",
		},
		Writer:      "file",
		OutputFile:  filepath.Join(t.TempDir(), "test_output.json"),
		FileMode:    "overwrite",
		PrettyPrint: true,
	}

	// Create filter
	msgFilter := filter.NewMessageFilter(cfg)
	if errs := msgFilter.GetCompilationErrors(); len(errs) > 0 {
		t.Fatalf("Filter compilation errors: %v", errs)
	}

	// Create exporter
	exp, err := exporter.NewExporter(cfg)
	if err != nil {
		t.Fatalf("Failed to create exporter: %v", err)
	}
	defer exp.Close()

	// Test message that should pass filter
	testMsg := &testMessage{body: []byte(`{"user": "admin", "action": "login"}`)}

	if !msgFilter.Filter(testMsg) {
		t.Error("Expected message to pass filter")
	}

	// Test message that should fail filter (contains debug)
	testMsg2 := &testMessage{body: []byte(`{"user": "admin", "action": "debug"}`)}

	if msgFilter.Filter(testMsg2) {
		t.Error("Expected message to fail filter (contains debug)")
	}

	// Test message that should fail filter (doesn't contain admin)
	testMsg3 := &testMessage{body: []byte(`{"user": "user", "action": "login"}`)}

	if msgFilter.Filter(testMsg3) {
		t.Error("Expected message to fail filter (doesn't contain admin)")
	}
}

func TestIntegration_JSONFilterAndExport(t *testing.T) {
	// Create a test configuration with JSON filter
	cfg := &config.Config{
		FilterConfig: struct {
			IncludePatterns []string
			ExcludePatterns []string
			JSONFilter      string
			MaxMessageSize  int
			RegexFilter     string
		}{
			IncludePatterns: []string{},
			ExcludePatterns: []string{},
			JSONFilter:      `.user.role == "admin"`,
			MaxMessageSize:  -1,
			RegexFilter:     "",
		},
		Writer:      "file",
		OutputFile:  filepath.Join(t.TempDir(), "test_output.json"),
		FileMode:    "overwrite",
		PrettyPrint: false,
	}

	// Create filter
	msgFilter := filter.NewMessageFilter(cfg)
	if errs := msgFilter.GetCompilationErrors(); len(errs) > 0 {
		t.Fatalf("Filter compilation errors: %v", errs)
	}

	// Create exporter
	exp, err := exporter.NewExporter(cfg)
	if err != nil {
		t.Fatalf("Failed to create exporter: %v", err)
	}
	defer exp.Close()

	// Test message that should pass JSON filter
	testMsg := &testMessage{body: []byte(`{"user": {"role": "admin", "name": "john"}}`)}

	if !msgFilter.Filter(testMsg) {
		t.Error("Expected message to pass JSON filter")
	}

	// Test message that should fail JSON filter
	testMsg2 := &testMessage{body: []byte(`{"user": {"role": "user", "name": "jane"}}`)}

	if msgFilter.Filter(testMsg2) {
		t.Error("Expected message to fail JSON filter")
	}
}

func TestIntegration_ConsoleExporter(t *testing.T) {
	// Create a test configuration for console export
	cfg := &config.Config{
		Writer:      "console",
		PrettyPrint: true,
	}

	// Create exporter
	exp, err := exporter.NewExporter(cfg)
	if err != nil {
		t.Fatalf("Failed to create console exporter: %v", err)
	}
	defer exp.Close()

	// Test that exporter implements the interface
	var _ exporter.Exporter = exp
}

func TestIntegration_FileExporterWithAppend(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test_append.json")

	// Create a test configuration for file export with append
	cfg := &config.Config{
		Writer:      "file",
		OutputFile:  tmpFile,
		FileMode:    "append",
		PrettyPrint: false,
	}

	// Create exporter
	exp, err := exporter.NewExporter(cfg)
	if err != nil {
		t.Fatalf("Failed to create file exporter: %v", err)
	}
	defer exp.Close()
	defer os.Remove(tmpFile)

	// Test that exporter implements the interface
	var _ exporter.Exporter = exp
}

func TestIntegration_ConfigCreation(t *testing.T) {
	// Test creating config with various options
	cfg := config.New(
		config.WithRabbitMQURL("amqp://localhost:5672/"),
		config.WithExchange("test_exchange"),
		config.WithQueue("test_queue"),
		config.WithRoutingKeys([]string{"user.created", "user.updated"}),
		config.WithAutoAck(true),
		config.WithSkipTLSVerification(true),
		config.WithOutputFile("test.json"),
		config.WithWriter("file"),
		config.WithPrettyPrint(true),
		config.WithIncludePatterns([]string{"admin"}),
		config.WithExcludePatterns([]string{"debug"}),
		config.WithJSONFilter(`.user.role == "admin"`),
		config.WithMaxMessageSize(1024),
	)

	if cfg == nil {
		t.Fatal("Expected config to be created")
	}

	// Verify all options were set
	if cfg.RabbitMQURL != "amqp://localhost:5672/" {
		t.Error("RabbitMQURL not set correctly")
	}

	if cfg.Exchange != "test_exchange" {
		t.Error("Exchange not set correctly")
	}

	if cfg.Queue != "test_queue" {
		t.Error("Queue not set correctly")
	}

	if len(cfg.RoutingKeys) != 2 {
		t.Error("RoutingKeys not set correctly")
	}

	if !cfg.AutoAck {
		t.Error("AutoAck not set correctly")
	}

	if !cfg.SkipTLSVerification {
		t.Error("SkipTLSVerification not set correctly")
	}

	if cfg.OutputFile != "test.json" {
		t.Error("OutputFile not set correctly")
	}

	if cfg.Writer != "file" {
		t.Error("Writer not set correctly")
	}

	if !cfg.PrettyPrint {
		t.Error("PrettyPrint not set correctly")
	}

	if len(cfg.FilterConfig.IncludePatterns) != 1 {
		t.Error("IncludePatterns not set correctly")
	}

	if len(cfg.FilterConfig.ExcludePatterns) != 1 {
		t.Error("ExcludePatterns not set correctly")
	}

	if cfg.FilterConfig.JSONFilter != `.user.role == "admin"` {
		t.Error("JSONFilter not set correctly")
	}

	if cfg.FilterConfig.MaxMessageSize != 1024 {
		t.Error("MaxMessageSize not set correctly")
	}
}

func TestIntegration_ErrorHandling(t *testing.T) {
	// Test invalid writer type
	cfg := &config.Config{
		Writer: "invalid_writer",
	}

	_, err := exporter.NewExporter(cfg)
	if err == nil {
		t.Error("Expected error for invalid writer type")
	}

	// Test invalid file mode
	cfg2 := &config.Config{
		Writer:     "file",
		OutputFile: filepath.Join(t.TempDir(), "test.json"),
		FileMode:   "invalid_mode",
	}

	_, err = exporter.NewFileWriter(cfg2)
	if err == nil {
		t.Error("Expected error for invalid file mode")
	}

	// Test invalid JSON filter
	cfg3 := &config.Config{
		FilterConfig: struct {
			IncludePatterns []string
			ExcludePatterns []string
			JSONFilter      string
			MaxMessageSize  int
			RegexFilter     string
		}{
			IncludePatterns: []string{},
			ExcludePatterns: []string{},
			JSONFilter:      "invalid jq syntax",
			MaxMessageSize:  -1,
			RegexFilter:     "",
		},
	}

	filter := filter.NewMessageFilter(cfg3)
	errs := filter.GetCompilationErrors()
	if len(errs) == 0 {
		t.Error("Expected compilation errors for invalid JSON filter")
	}
}

// testMessage implements filter.MessageDelivery interface for testing
type testMessage struct {
	body []byte
}

func (t *testMessage) GetBody() []byte {
	return t.body
}
