package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	RabbitMQURL         string
	Exchange            string
	Queue               string
	OutputFile          string
	UseAMQPS            bool
	VirtualHost         string
	SkipTLSVerification bool
	AutoAck             bool
	FileMode            string
	StopAfterConsume    bool
	RoutingKeys         []string
	PrettyPrint         bool
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

func New(options ...Option) *Config {
	c := &Config{
		RabbitMQURL:         fmt.Sprintf("%s://%s/%s", getProtocol(), viper.GetString("url"), viper.GetString("virtualhost")),
		Exchange:            viper.GetString("exchange"),
		Queue:               viper.GetString("queue"),
		OutputFile:          viper.GetString("output"),
		UseAMQPS:            viper.GetBool("amqps"),
		VirtualHost:         viper.GetString("virtualhost"),
		SkipTLSVerification: viper.GetBool("skip-tls-verify"),
		AutoAck:             viper.GetBool("auto-ack"),
		FileMode:            viper.GetString("file-mode"),
		StopAfterConsume:    viper.GetBool("stop-after-consume"),
		RoutingKeys:         viper.GetStringSlice("routing-keys"),
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *Config) PrintConfig() (string, error) {
	configJSON, err := json.MarshalIndent(c, "", "  ")
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
