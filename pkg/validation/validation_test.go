package validation

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestValidateInput_EmptyURL(t *testing.T) {
	resetViper()
	// Set empty URL
	viper.Set("url", "")

	err := ValidateInput()
	if err == nil {
		t.Error("Expected error for empty URL")
	}

	if err.Error() != "RabbitMQ URL is required" {
		t.Errorf("Expected 'RabbitMQ URL is required', got: %s", err.Error())
	}
}

func TestValidateInput_ValidURL(t *testing.T) {
	resetViper()
	// Set valid URL
	viper.Set("url", "localhost:5672")

	err := ValidateInput()
	if err != nil {
		t.Errorf("Unexpected error for valid URL: %v", err)
	}
}

func TestValidateInput_InvalidURL(t *testing.T) {
	resetViper()
	// Set invalid URL that will definitely fail parsing
	viper.Set("url", "[invalid")

	err := ValidateInput()
	if err == nil {
		t.Error("Expected error for invalid URL")
		return
	}

	if !strings.Contains(err.Error(), "invalid RabbitMQ URL") {
		t.Errorf("Expected URL parse error, got: %s", err.Error())
	}
}

func TestValidateInput_SecureURL(t *testing.T) {
	resetViper()
	// Set secure URL
	viper.Set("url", "localhost:5671")
	viper.Set("secure", true)

	err := ValidateInput()
	if err != nil {
		t.Errorf("Unexpected error for secure URL: %v", err)
	}
}

func TestValidateVirtualHost_Default(t *testing.T) {
	// Test default virtual host
	viper.Set("virtualhost", "/")

	err := validateVirtualHost()
	if err != nil {
		t.Errorf("Unexpected error for default virtual host: %v", err)
	}
}

func TestValidateVirtualHost_Empty(t *testing.T) {
	// Test empty virtual host
	viper.Set("virtualhost", "")

	err := validateVirtualHost()
	if err != nil {
		t.Errorf("Unexpected error for empty virtual host: %v", err)
	}
}

func TestValidateVirtualHost_Valid(t *testing.T) {
	// Test valid virtual host
	viper.Set("virtualhost", "my-vhost")

	err := validateVirtualHost()
	if err != nil {
		t.Errorf("Unexpected error for valid virtual host: %v", err)
	}
}

func TestValidateVirtualHost_Invalid(t *testing.T) {
	// Test invalid virtual host (starts with /)
	viper.Set("virtualhost", "/invalid")

	err := validateVirtualHost()
	if err == nil {
		t.Error("Expected error for invalid virtual host")
	}

	if err.Error() != "virtual host should not start with '/' (except for the default '/' virtual host)" {
		t.Errorf("Expected virtual host error, got: %s", err.Error())
	}
}

func TestValidateWriter_Console(t *testing.T) {
	resetViper()
	// Test console writer
	viper.Set("writer", "console")

	err := ValidateInput()
	if err != nil {
		t.Errorf("Unexpected error for console writer: %v", err)
	}
}

func TestValidateWriter_File(t *testing.T) {
	resetViper()
	// Test file writer
	viper.Set("writer", "file")
	viper.Set("output", "test.json")

	err := ValidateInput()
	if err != nil {
		t.Errorf("Unexpected error for file writer: %v", err)
	}
}

func TestValidateWriter_Invalid(t *testing.T) {
	resetViper()
	// Test invalid writer
	viper.Set("writer", "invalid")

	err := ValidateInput()
	if err == nil {
		t.Error("Expected error for invalid writer")
	}

	if !strings.Contains(err.Error(), "invalid writer type") {
		t.Errorf("Expected writer error, got: %s", err.Error())
	}
}

func TestValidateWriter_FileWithoutOutput(t *testing.T) {
	resetViper()
	// Test file writer without output file
	viper.Set("writer", "file")
	viper.Set("output", "")

	err := ValidateInput()
	if err == nil {
		t.Error("Expected error for file writer without output")
	}

	if !strings.Contains(err.Error(), "output file is required") {
		t.Errorf("Expected output file error, got: %s", err.Error())
	}
}

func TestValidateWriter_InvalidFileMode(t *testing.T) {
	resetViper()
	// Test invalid file mode
	viper.Set("writer", "file")
	viper.Set("output", "test.json")
	viper.Set("file-mode", "invalid")

	err := ValidateInput()
	if err == nil {
		t.Error("Expected error for invalid file mode")
	}

	if !strings.Contains(err.Error(), "invalid file mode") {
		t.Errorf("Expected file mode error, got: %s", err.Error())
	}
}

func TestValidatePatterns_EmptyInclude(t *testing.T) {
	resetViper()
	// Test empty include pattern
	viper.Set("include-patterns", []string{"", "valid"})

	err := ValidateInput()
	if err == nil {
		t.Error("Expected error for empty include pattern")
	}

	if !strings.Contains(err.Error(), "include patterns cannot be empty") {
		t.Errorf("Expected include pattern error, got: %s", err.Error())
	}
}

func TestValidatePatterns_EmptyExclude(t *testing.T) {
	resetViper()
	// Test empty exclude pattern
	viper.Set("exclude-patterns", []string{"", "valid"})

	err := ValidateInput()
	if err == nil {
		t.Error("Expected error for empty exclude pattern")
	}

	if !strings.Contains(err.Error(), "exclude patterns cannot be empty") {
		t.Errorf("Expected exclude pattern error, got: %s", err.Error())
	}
}

func TestValidateRegex_Valid(t *testing.T) {
	resetViper()
	// Test valid regex
	viper.Set("regex-filter", "error.*")

	err := ValidateInput()
	if err != nil {
		t.Errorf("Unexpected error for valid regex: %v", err)
	}
}

func TestValidateRegex_Invalid(t *testing.T) {
	resetViper()
	// Test invalid regex
	viper.Set("regex-filter", "[invalid")

	err := ValidateInput()
	if err == nil {
		t.Error("Expected error for invalid regex")
	}

	if !strings.Contains(err.Error(), "invalid regex filter") {
		t.Errorf("Expected regex error, got: %s", err.Error())
	}
}

func TestValidateMessageSize_Valid(t *testing.T) {
	resetViper()
	// Test valid max message size
	viper.Set("max-message-size", 1024)

	err := ValidateInput()
	if err != nil {
		t.Errorf("Unexpected error for valid max message size: %v", err)
	}
}

func TestValidateMessageSize_Unlimited(t *testing.T) {
	resetViper()
	// Test unlimited max message size
	viper.Set("max-message-size", -1)

	err := ValidateInput()
	if err != nil {
		t.Errorf("Unexpected error for unlimited max message size: %v", err)
	}
}

func TestValidateMessageSize_Invalid(t *testing.T) {
	resetViper()
	// Test invalid max message size
	viper.Set("max-message-size", -5)

	err := ValidateInput()
	if err == nil {
		t.Error("Expected error for invalid max message size")
	}

	if !strings.Contains(err.Error(), "max message size must be -1 or a positive integer") {
		t.Errorf("Expected max message size error, got: %s", err.Error())
	}
}

func TestValidateMessageSize_Zero(t *testing.T) {
	resetViper()
	// Test zero max message size
	viper.Set("max-message-size", 0)

	err := ValidateInput()
	if err == nil {
		t.Error("Expected error for zero max message size")
	}

	if !strings.Contains(err.Error(), "max message size must be -1 or a positive integer") {
		t.Errorf("Expected max message size error, got: %s", err.Error())
	}
}

func TestValidateInput_CompleteValidConfig(t *testing.T) {
	resetViper()
	// Test complete valid configuration
	viper.Set("url", "localhost:5672")
	viper.Set("virtualhost", "/")
	viper.Set("writer", "console")
	viper.Set("file-mode", "overwrite")
	viper.Set("max-message-size", 1024)

	err := ValidateInput()
	if err != nil {
		t.Errorf("Unexpected error for complete valid config: %v", err)
	}
}

func TestValidateInput_FileWriterWithOutput(t *testing.T) {
	resetViper()
	// Test file writer with output file
	viper.Set("url", "localhost:5672")
	viper.Set("writer", "file")
	viper.Set("output", "test.json")
	viper.Set("file-mode", "append")

	err := ValidateInput()
	if err != nil {
		t.Errorf("Unexpected error for file writer with output: %v", err)
	}
}

// Helper function to reset viper for each test
func resetViper() {
	viper.Reset()
	// Set some defaults to avoid issues
	viper.Set("url", "localhost:5672")
	viper.Set("virtualhost", "/")
	viper.Set("writer", "console")
	viper.Set("file-mode", "overwrite")
	viper.Set("max-message-size", -1)
}

func init() {
	// Reset viper before each test
	resetViper()
}
