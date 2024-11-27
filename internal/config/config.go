package config

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/itchyny/gojq"
	"github.com/spf13/viper"
)

type ExporterKind string

const (
	ConsoleExporterKind ExporterKind = "console"
	FileWriterKind      ExporterKind = "file"
)

type Config struct {
	RabbitMQURL         string
	Exchange            string
	Queue               string
	Writer              ExporterKind
	OutputFile          string
	FileMode            string
	UseAMQPS            bool
	VirtualHost         string
	SkipTLSVerification bool
	AutoAck             bool
	StopAfterConsume    bool
	RoutingKeys         []string
	PrettyPrint         bool
	FullMessage         bool

	FilterConfig struct {
		IncludePatterns []string
		ExcludePatterns []string
		JSONFilter      *gojq.Query
		MaxMessageSize  int
		RegexPattern    string
		CompileRegex    *regexp.Regexp
	}
}

type Option func(*Config)

func WithRabbitMQURL(url string) Option {
	return func(c *Config) {
		c.RabbitMQURL = url
	}
}

func WithExchange(exchange string) Option {
	return func(c *Config) {
		c.Exchange = exchange
	}
}

func WithQueue(queue string) Option {
	return func(c *Config) {
		c.Queue = queue
	}
}

func WithOutputFile(outputFile string) Option {
	return func(c *Config) {
		c.OutputFile = outputFile
	}
}

func WithUseAMQPS(useAMQPS bool) Option {
	return func(c *Config) {
		c.UseAMQPS = useAMQPS
	}
}

func WithVirtualHost(virtualHost string) Option {
	return func(c *Config) {
		c.VirtualHost = virtualHost
	}
}

func WithSkipTLSVerification(skip bool) Option {
	return func(c *Config) {
		c.SkipTLSVerification = skip
	}
}

func WithAutoAck(autoAck bool) Option {
	return func(c *Config) {
		c.AutoAck = autoAck
	}
}

func WithFileMode(fileMode string) Option {
	return func(c *Config) {
		c.FileMode = fileMode
	}
}

func WithStopAfterConsume(stop bool) Option {
	return func(c *Config) {
		c.StopAfterConsume = stop
	}
}

func WithRoutingKeys(routingKeys []string) Option {
	return func(c *Config) {
		c.RoutingKeys = routingKeys
	}
}

func WithPrettyPrint(prettyPrint bool) Option {
	return func(c *Config) {
		c.PrettyPrint = prettyPrint
	}
}

// Add new Option functions
func WithIncludePatterns(patterns []string) Option {
	return func(c *Config) {
		c.FilterConfig.IncludePatterns = patterns
	}
}

func WithExcludePatterns(patterns []string) Option {
	return func(c *Config) {
		c.FilterConfig.ExcludePatterns = patterns
	}
}

// Modify the WithJSONFilter function
func WithJSONFilter(jsonFilter string) Option {
	return func(c *Config) {
		if jsonFilter == "" {
			return
		}
		// Parse the query during configuration
		query, err := gojq.Parse(jsonFilter)
		if err != nil {
			// Handle error - you might want to log or handle this differently
			fmt.Printf("Error parsing JSON filter: %v\n", err)
			return
		}
		c.FilterConfig.JSONFilter = query
	}
}

func WithMaxMessageSize(size int) Option {
	return func(c *Config) {
		c.FilterConfig.MaxMessageSize = size
	}
}

func WithRegexFilter(pattern string) Option {
	return func(c *Config) {
		c.FilterConfig.RegexPattern = pattern
		if regex, err := regexp.Compile(pattern); err == nil {
			c.FilterConfig.CompileRegex = regex
		}
	}
}

func WithFullMessage(fullMessage bool) Option {
	return func(c *Config) {
		c.FullMessage = fullMessage
	}
}

func WithWriter(writer string) Option {
	return func(c *Config) {
		c.Writer = ExporterKind(writer)
	}
}

func New(options ...Option) *Config {
	c := &Config{
		RabbitMQURL:         fmt.Sprintf("%s://%s/%s", getProtocol(), viper.GetString("url"), viper.GetString("virtualhost")),
		Exchange:            viper.GetString("exchange"),
		Queue:               viper.GetString("queue"),
		Writer:              ExporterKind(viper.GetString("writer")),
		OutputFile:          viper.GetString("output"),
		FileMode:            viper.GetString("file-mode"),
		UseAMQPS:            viper.GetBool("amqps"),
		VirtualHost:         viper.GetString("virtualhost"),
		SkipTLSVerification: viper.GetBool("skip-tls-verify"),
		AutoAck:             viper.GetBool("auto-ack"),
		StopAfterConsume:    viper.GetBool("stop-after-consume"),
		RoutingKeys:         viper.GetStringSlice("routing-keys"),
		PrettyPrint:         viper.GetBool("pretty-print"),
		FullMessage:         viper.GetBool("full-message"),

		FilterConfig: struct {
			IncludePatterns []string
			ExcludePatterns []string
			JSONFilter      *gojq.Query
			MaxMessageSize  int
			RegexPattern    string
			CompileRegex    *regexp.Regexp
		}{
			MaxMessageSize: -1, // Default: no size limit
		},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *Config) PrintConfig() (string, error) {
	configJSON, err := json.MarshalIndent(map[string]interface{}{
		"RabbitMQURL":         c.RabbitMQURL,
		"Exchange":            c.Exchange,
		"Queue":               c.Queue,
		"Writer":              c.Writer,
		"OutputFile":          c.OutputFile,
		"FileMode":            c.FileMode,
		"UseAMQPS":            c.UseAMQPS,
		"VirtualHost":         c.VirtualHost,
		"SkipTLSVerification": c.SkipTLSVerification,
		"AutoAck":             c.AutoAck,
		"StopAfterConsume":    c.StopAfterConsume,
		"RoutingKeys":         c.RoutingKeys,
		"PrettyPrint":         c.PrettyPrint,
		"FullMessage":         c.FullMessage,
		"FilterConfig": map[string]interface{}{
			"IncludePatterns": c.FilterConfig.IncludePatterns,
			"ExcludePatterns": c.FilterConfig.ExcludePatterns,
			"JSONFilter": func() string {
				if c.FilterConfig.JSONFilter != nil {
					return c.FilterConfig.JSONFilter.String()
				}
				return ""
			}(),
			"MaxMessageSize": c.FilterConfig.MaxMessageSize,
			"RegexPattern":   c.FilterConfig.RegexPattern,
		},
	}, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %v", err)
	}

	return string(configJSON), nil
}

func getProtocol() string {
	if viper.GetBool("amqps") {
		return "amqps"
	}
	return "amqp"
}
