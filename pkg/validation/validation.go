package validation

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

var (
	ValidWriters   = []string{"file", "console"}
	ValidFileModes = []string{"append", "overwrite"}
)

func ValidateInput() error {
	if err := validateURL(); err != nil {
		return err
	}
	if err := validateVirtualHost(); err != nil {
		return err
	}
	if err := validateWriter(); err != nil {
		return err
	}
	if err := validatePatterns(); err != nil {
		return err
	}
	if err := validateRegex(); err != nil {
		return err
	}
	return validateMessageSize()
}

func validateURL() error {
	urlStr := viper.GetString("url")
	if urlStr == "" {
		return fmt.Errorf("RabbitMQ URL is required")
	}
	protocol := "amqp"
	if viper.GetBool("secure") {
		protocol = "amqps"
	}
	if _, err := url.Parse(fmt.Sprintf("%s://%s", protocol, urlStr)); err != nil {
		return fmt.Errorf("invalid RabbitMQ URL: %v", err)
	}
	return nil
}

func validateVirtualHost() error {
	vh := viper.GetString("virtualhost")
	// Allow "/" as it's the default virtual host in RabbitMQ
	if vh != "" && vh != "/" && strings.HasPrefix(vh, "/") {
		return fmt.Errorf("virtual host should not start with '/' (except for the default '/' virtual host)")
	}
	return nil
}

func validateWriter() error {
	writer := viper.GetString("writer")
	if !contains(ValidWriters, writer) {
		return fmt.Errorf("invalid writer type '%s', must be one of: %v", writer, ValidWriters)
	}
	if writer == "file" && viper.GetString("output") == "" {
		return fmt.Errorf("output file is required when using file writer")
	}
	if !contains(ValidFileModes, viper.GetString("file-mode")) {
		return fmt.Errorf("invalid file mode '%s': must be one of: %v", viper.GetString("file-mode"), ValidFileModes)
	}
	return nil
}

func validatePatterns() error {
	for _, pattern := range viper.GetStringSlice("include-patterns") {
		if strings.TrimSpace(pattern) == "" {
			return fmt.Errorf("include patterns cannot be empty")
		}
	}
	for _, pattern := range viper.GetStringSlice("exclude-patterns") {
		if strings.TrimSpace(pattern) == "" {
			return fmt.Errorf("exclude patterns cannot be empty")
		}
	}
	return nil
}

func validateRegex() error {
	if regexFilter := viper.GetString("regex-filter"); regexFilter != "" {
		if _, err := regexp.Compile(regexFilter); err != nil {
			return fmt.Errorf("invalid regex filter: %v", err)
		}
	}
	return nil
}

func validateMessageSize() error {
	if size := viper.GetInt("max-message-size"); size != -1 && size <= 0 {
		return fmt.Errorf("max message size must be -1 or a positive integer")
	}
	return nil
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
