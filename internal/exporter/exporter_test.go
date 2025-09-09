package exporter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marianozunino/goq/internal/config"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
)

func TestNewConsoleExporter(t *testing.T) {
	cfg := &config.Config{}
	exporter, err := NewConsoleExporter(cfg)

	if err != nil {
		t.Fatalf("Unexpected error creating console exporter: %v", err)
	}

	if exporter == nil {
		t.Fatal("Expected exporter to be created")
	}

	// Console exporter should implement Exporter interface
	var _ Exporter = exporter
}

func TestConsoleExporter_WriteMessage(t *testing.T) {
	cfg := &config.Config{PrettyPrint: true}
	exporter, err := NewConsoleExporter(cfg)
	if err != nil {
		t.Fatalf("Failed to create console exporter: %v", err)
	}

	// Create a test message
	var msg rabbitmq.Delivery
	msg.Body = []byte(`{"test": "data"}`)
	msg.Exchange = "test_exchange"
	msg.RoutingKey = "test.key"
	msg.Headers = amqp091.Table{"test-header": "test-value"}

	// This should not panic or error
	err = exporter.WriteMessage(msg)
	if err != nil {
		t.Errorf("Unexpected error writing message: %v", err)
	}
}

func TestConsoleExporter_Close(t *testing.T) {
	cfg := &config.Config{}
	exporter, err := NewConsoleExporter(cfg)
	if err != nil {
		t.Fatalf("Failed to create console exporter: %v", err)
	}

	// Close should not error
	err = exporter.Close()
	if err != nil {
		t.Errorf("Unexpected error closing console exporter: %v", err)
	}
}

func TestNewFileWriter_AppendMode(t *testing.T) {
	// Create a temporary file
	tmpFile := filepath.Join(t.TempDir(), "test_output.json")

	cfg := &config.Config{
		OutputFile: tmpFile,
		FileMode:   "append",
	}

	exporter, err := NewFileWriter(cfg)
	if err != nil {
		t.Fatalf("Unexpected error creating file writer: %v", err)
	}

	if exporter == nil {
		t.Fatal("Expected exporter to be created")
	}

	// Clean up
	exporter.Close()
	os.Remove(tmpFile)
}

func TestNewFileWriter_OverwriteMode(t *testing.T) {
	// Create a temporary file
	tmpFile := filepath.Join(t.TempDir(), "test_output.json")

	cfg := &config.Config{
		OutputFile: tmpFile,
		FileMode:   "overwrite",
	}

	exporter, err := NewFileWriter(cfg)
	if err != nil {
		t.Fatalf("Unexpected error creating file writer: %v", err)
	}

	if exporter == nil {
		t.Fatal("Expected exporter to be created")
	}

	// Clean up
	exporter.Close()
	os.Remove(tmpFile)
}

func TestNewFileWriter_InvalidMode(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test_output.json")

	cfg := &config.Config{
		OutputFile: tmpFile,
		FileMode:   "invalid_mode",
	}

	exporter, err := NewFileWriter(cfg)
	if err == nil {
		t.Error("Expected error for invalid file mode")
	}

	if exporter != nil {
		t.Error("Expected exporter to be nil for invalid mode")
	}
}

func TestFileExporter_WriteMessage(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test_output.json")

	cfg := &config.Config{
		OutputFile:  tmpFile,
		FileMode:    "overwrite",
		PrettyPrint: true,
	}

	exporter, err := NewFileWriter(cfg)
	if err != nil {
		t.Fatalf("Failed to create file writer: %v", err)
	}
	defer exporter.Close()
	defer os.Remove(tmpFile)

	// Create a test message
	var msg rabbitmq.Delivery
	msg.Body = []byte(`{"test": "data"}`)
	msg.Exchange = "test_exchange"
	msg.RoutingKey = "test.key"
	msg.Headers = amqp091.Table{"test-header": "test-value"}

	err = exporter.WriteMessage(msg)
	if err != nil {
		t.Errorf("Unexpected error writing message: %v", err)
	}

	// Check if file was created and has content
	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Errorf("File was not created: %v", err)
	}

	if info.Size() == 0 {
		t.Error("File is empty after writing message")
	}
}

func TestFileExporter_WriteMultipleMessages(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test_output.json")

	cfg := &config.Config{
		OutputFile:  tmpFile,
		FileMode:    "overwrite",
		PrettyPrint: false,
	}

	exporter, err := NewFileWriter(cfg)
	if err != nil {
		t.Fatalf("Failed to create file writer: %v", err)
	}
	defer exporter.Close()
	defer os.Remove(tmpFile)

	// Write multiple messages
	for i := 0; i < 3; i++ {
		var msg rabbitmq.Delivery
		msg.Body = []byte(`{"message": ` + string(rune(i+48)) + `}`)
		msg.Exchange = "test_exchange"
		msg.RoutingKey = "test.key"
		msg.Headers = amqp091.Table{"index": i}

		err = exporter.WriteMessage(msg)
		if err != nil {
			t.Errorf("Unexpected error writing message %d: %v", i, err)
		}
	}

	// Check file content
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Errorf("Failed to read output file: %v", err)
	}

	// Should contain multiple JSON objects
	if len(content) == 0 {
		t.Error("File is empty after writing multiple messages")
	}
}

func TestNewExporter_ConsoleWriter(t *testing.T) {
	cfg := &config.Config{Writer: config.ConsoleExporterKind}

	exporter, err := NewExporter(cfg)
	if err != nil {
		t.Fatalf("Unexpected error creating console exporter: %v", err)
	}

	if exporter == nil {
		t.Fatal("Expected exporter to be created")
	}
}

func TestNewExporter_FileWriter(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test_output.json")

	cfg := &config.Config{
		Writer:     config.FileWriterKind,
		OutputFile: tmpFile,
		FileMode:   "overwrite",
	}

	exporter, err := NewExporter(cfg)
	if err != nil {
		t.Fatalf("Unexpected error creating file exporter: %v", err)
	}

	if exporter == nil {
		t.Fatal("Expected exporter to be created")
	}

	exporter.Close()
	os.Remove(tmpFile)
}

func TestNewExporter_InvalidWriter(t *testing.T) {
	cfg := &config.Config{Writer: "invalid_writer"}

	exporter, err := NewExporter(cfg)
	if err == nil {
		t.Error("Expected error for invalid writer type")
	}

	if exporter != nil {
		t.Error("Expected exporter to be nil for invalid writer")
	}
}
