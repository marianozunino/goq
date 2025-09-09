package exporter

import (
	"fmt"
	"testing"

	"github.com/marianozunino/goq/internal/config"
)

func TestExporterError_Error(t *testing.T) {
	err := &ExporterError{
		Type: "file",
		Err:  fmt.Errorf("test error"),
	}

	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("Expected Error() to return non-empty string")
	}

	if !contains(errorMsg, "test error") {
		t.Error("Expected error message to contain 'test error'")
	}

	if !contains(errorMsg, "file") {
		t.Error("Expected error message to contain exporter type")
	}
}

func TestExporterError_Unwrap(t *testing.T) {
	originalErr := fmt.Errorf("wrapped error")

	wrappedErr := &ExporterError{
		Type: "file",
		Err:  originalErr,
	}

	unwrapped := wrappedErr.Unwrap()
	if unwrapped != originalErr {
		t.Error("Expected Unwrap() to return the original error")
	}
}

func TestConvertHeaders_ComplexTypes(t *testing.T) {
	// Test with complex header types
	headers := map[string]interface{}{
		"string-header": "test-value",
		"int-header":    42,
		"float-header":  3.14,
		"bool-header":   true,
		"array-header":  []interface{}{"item1", "item2"},
		"map-header": map[string]interface{}{
			"nested-key": "nested-value",
		},
	}

	converted := convertHeaders(headers)

	if converted == nil {
		t.Error("Expected converted headers to not be nil")
	}

	// Verify that the conversion preserves the structure
	if len(converted) != len(headers) {
		t.Errorf("Expected %d headers, got %d", len(headers), len(converted))
	}
}

func TestConvertHeaders_Nil(t *testing.T) {
	converted := convertHeaders(nil)

	if converted != nil {
		t.Error("Expected nil headers to remain nil")
	}
}

func TestConvertHeaders_Empty(t *testing.T) {
	headers := map[string]interface{}{}

	converted := convertHeaders(headers)

	if converted == nil {
		t.Error("Expected converted headers to not be nil")
	}

	if len(converted) != 0 {
		t.Error("Expected empty headers to remain empty")
	}
}

func TestFileExporter_ErrorHandling(t *testing.T) {
	// Test file exporter with invalid configuration
	cfg := &config.Config{
		Writer:     "file",
		OutputFile: "", // Empty output file
		FileMode:   "overwrite",
	}

	exporter, err := NewFileWriter(cfg)
	if err == nil {
		t.Error("Expected error for empty output file")
	}

	if exporter != nil {
		t.Error("Expected exporter to be nil for empty output file")
	}
}

func TestConsoleExporter_ErrorHandling(t *testing.T) {
	// Console exporter should not fail with any config
	cfg := &config.Config{
		Writer:      "console",
		PrettyPrint: true,
	}

	exporter, err := NewConsoleExporter(cfg)
	if err != nil {
		t.Errorf("Unexpected error creating console exporter: %v", err)
	}

	if exporter == nil {
		t.Error("Expected exporter to be created")
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
