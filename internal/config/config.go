package config

import (
	"fmt"
	"strings"

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
		JSONFilter      string
		MaxMessageSize  int
		RegexFilter     string
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
		c.FilterConfig.JSONFilter = jsonFilter
	}
}

func WithMaxMessageSize(size int) Option {
	return func(c *Config) {
		c.FilterConfig.MaxMessageSize = size
	}
}

func WithRegexFilter(pattern string) Option {
	return func(c *Config) {
		c.FilterConfig.RegexFilter = pattern
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
			JSONFilter      string
			MaxMessageSize  int
			RegexFilter     string
		}{
			MaxMessageSize: -1, // Default: no size limit
		},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *Config) PrintConfig() string {
	return fmt.Sprintf(`RabbitMQ:
	URL: %s
	Exchange: %s
	Queue: %s
	Virtual Host: %s
	Skip TLS Verification: %v
	Auto Acknowledge: %v
	Stop After Consume: %v
	Routing Keys: %v
Writer:
	Writer: %s
	Output File: %s
	File Mode: %s
	Pretty Print: %v
	Full Message: %v
Filters:
	Include Patterns: %s
	Exclude Patterns: %s
	JSON Filter: %s
	Max Message Size: %s
	Regex Filter: %s`,
		// RabbitMQ Section
		c.RabbitMQURL,
		c.Exchange,
		c.Queue,
		c.VirtualHost,
		c.SkipTLSVerification,
		c.AutoAck,
		c.StopAfterConsume,
		func() string {
			if len(c.RoutingKeys) == 0 {
				return "false"
			}
			return strings.Join(c.RoutingKeys, ", ")
		}(),
		// Writer Section
		c.Writer,
		func() string {
			if c.Writer == FileWriterKind {
				return c.OutputFile
			}
			return "false"
		}(),
		c.FileMode,
		c.PrettyPrint,
		c.FullMessage,
		// Filter Section
		func() string {
			if len(c.FilterConfig.IncludePatterns) == 0 {
				return "false"
			}
			return strings.Join(c.FilterConfig.IncludePatterns, ", ")
		}(),
		func() string {
			if len(c.FilterConfig.ExcludePatterns) == 0 {
				return "false"
			}
			return strings.Join(c.FilterConfig.ExcludePatterns, ", ")
		}(),
		func() string {
			if c.FilterConfig.JSONFilter == "" {
				return "false"
			}
			return c.FilterConfig.JSONFilter
		}(),
		func() string {
			if c.FilterConfig.MaxMessageSize == -1 {
				return "false"
			}
			return fmt.Sprintf("%d bytes", c.FilterConfig.MaxMessageSize)
		}(),
		func() string {
			if c.FilterConfig.RegexFilter == "" {
				return "false"
			}
			return c.FilterConfig.RegexFilter
		}(),
	)
}

func getProtocol() string {
	if viper.GetBool("amqps") {
		return "amqps"
	}
	return "amqp"
}
