package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/adrg/xdg"
	"github.com/marianozunino/goq/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultURL         = "localhost:5672"
	defaultVirtualHost = "/"
)

func InitConfig() {
	configPath := viper.GetString("config")
	if configPath == "" {
		var err error
		configPath, err = xdg.ConfigFile("goq/goq.yaml")
		if err != nil {
			log.Fatal(err)
		}
	}
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.ReadInConfig()
}

func SetupFlags(flags *pflag.FlagSet, validWriters, validFileModes []string) {
	flags.StringP("exchange", "e", "", "RabbitMQ exchange name")
	flags.StringP("url", "u", defaultURL, "RabbitMQ server URL")
	flags.StringP("virtualhost", "v", defaultVirtualHost, "RabbitMQ virtual host")
	flags.BoolP("amqps", "s", false, "Use AMQPS instead of AMQP")
	flags.BoolP("tls-skip", "k", false, "Skip TLS certificate verification")
	flags.StringP("writer", "w", "file", fmt.Sprintf("Output writer type (%s)", strings.Join(validWriters, " or ")))
	flags.StringP("output", "o", "", "Output file name")
	flags.StringP("file-mode", "m", "overwrite", fmt.Sprintf("File mode (%s)", strings.Join(validFileModes, " or ")))
	flags.BoolP("pretty-print", "p", false, "Pretty print JSON messages")
	flags.BoolP("full-message", "f", false, "Print complete message details")
	flags.StringSliceP("include-patterns", "i", []string{}, "Include messages containing these patterns")
	flags.StringSliceP("exclude-patterns", "x", []string{}, "Exclude messages containing these patterns")
	flags.StringP("json-filter", "j", "", "JSON filter expression")
	flags.StringP("regex-filter", "R", "", "Regex pattern to filter messages")
	flags.IntP("max-message-size", "z", -1, "Maximum message size in bytes (-1 for unlimited)")
	flags.String("config", xdg.ConfigHome+"/goq/goq.yaml", "Config file path")

	viper.BindPFlags(flags)
}

func CreateCommonConfig(cmd *cobra.Command) *config.Config {
	queue, _ := cmd.Flags().GetString("queue")
	routingKeys, _ := cmd.Flags().GetStringSlice("routing-keys")
	autoAck, _ := cmd.Flags().GetBool("auto-ack")
	stopAfterConsume, _ := cmd.Flags().GetBool("stop-after-consume")

	protocol := "amqp"
	if viper.GetBool("amqps") {
		protocol = "amqps"
	}

	options := []config.Option{
		config.WithRabbitMQURL(fmt.Sprintf("%s://%s/%s", protocol, viper.GetString("url"), viper.GetString("virtualhost"))),
		config.WithExchange(viper.GetString("exchange")),
		config.WithVirtualHost(viper.GetString("virtualhost")),
		config.WithSkipTLSVerification(viper.GetBool("tls-skip")),
		config.WithQueue(queue),
		config.WithRoutingKeys(routingKeys),
		config.WithAutoAck(autoAck),
		config.WithStopAfterConsume(stopAfterConsume),
		config.WithOutputFile(viper.GetString("output")),
		config.WithFileMode(viper.GetString("file-mode")),
		config.WithWriter(viper.GetString("writer")),
		config.WithPrettyPrint(viper.GetBool("pretty-print")),
		config.WithQueue(queue),
		config.WithStopAfterConsume(stopAfterConsume),
		config.WithAutoAck(autoAck),
		config.WithFullMessage(viper.GetBool("full-message")),
		config.WithIncludePatterns(viper.GetStringSlice("include-patterns")),
		config.WithExcludePatterns(viper.GetStringSlice("exclude-patterns")),
		config.WithMaxMessageSize(viper.GetInt("max-message-size")),
		config.WithRegexFilter(viper.GetString("regex-filter")),
		config.WithJSONFilter(viper.GetString("json-filter")),
	}

	return config.New(options...)
}
