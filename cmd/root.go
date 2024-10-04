/*
Copyright © 2024 Mariano Zunino <marianoz@posteo.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/adrg/xdg"
	"github.com/fatih/color"
	app "github.com/marianozunino/goq/internal"
	"github.com/marianozunino/goq/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var logo = `
  ______    ______    ______
 /\  ___\  /\  __ \  /\  __ \
 \ \ \__ \ \ \ \/\ \ \ \ \/\_\
  \ \_____\ \ \_____\ \ \___\_\
   \/_____/  \/_____/  \/___/_/ ` + VersionFromBuild()

var rootCmd = &cobra.Command{
	Use:   "goq",
	Short: "A tool to dump RabbitMQ messages to a file",
	Long: logo + `

This application connects to a RabbitMQ server, consumes messages from a specified queue,
and writes them to a file while keeping the messages in the queue.`,
	Run: run,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/goq/goq.yaml)")
	rootCmd.Flags().StringP("url", "u", "", "RabbitMQ URL (e.g., localhost:5672)")
	rootCmd.Flags().StringP("exchange", "e", "", "RabbitMQ exchange name")
	rootCmd.Flags().StringP("queue", "q", "", "RabbitMQ queue name")
	rootCmd.Flags().StringP("output", "o", "messages.txt", "Output file name")
	rootCmd.Flags().BoolP("amqps", "s", false, "Use AMQPS instead of AMQP")
	rootCmd.Flags().StringP("virtualhost", "v", "", "RabbitMQ virtual host")
	rootCmd.Flags().BoolP("skip-tls-verify", "k", false, "Skip TLS certificate verification (insecure)")
	rootCmd.Flags().BoolP("auto-ack", "a", false, "Automatically acknowledge messages")
	rootCmd.Flags().StringP("file-mode", "m", "overwrite", "File mode (append or overwrite)")
	rootCmd.Flags().BoolP("stop-after-consume", "c", false, "Stop consuming after getting all messages from the queue")

	viper.BindPFlags(rootCmd.Flags())
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		configPath, err := xdg.ConfigFile("goq/goq.yaml")
		if err != nil {
			log.Fatal(err)
		}
		viper.SetConfigFile(configPath)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func run(cmd *cobra.Command, args []string) {
	// Validate required fields
	if err := validateRequiredFields(); err != nil {
		color.Red("Validation error: %v", err)
		os.Exit(1)
	}

	cfg := config.New(
		config.WithRabbitMQURL(fmt.Sprintf("%s://%s/%s", getProtocol(), viper.GetString("url"), viper.GetString("virtualhost"))),
		config.WithExchange(viper.GetString("exchange")),
		config.WithQueue(viper.GetString("queue")),
		config.WithOutputFile(viper.GetString("output")),
		config.WithUseAMQPS(viper.GetBool("amqps")),
		config.WithVirtualHost(viper.GetString("virtualhost")),
		config.WithSkipTLSVerification(viper.GetBool("skip-tls-verify")),
		config.WithAutoAck(viper.GetBool("auto-ack")),
		config.WithFileMode(viper.GetString("file-mode")),
		config.WithStopAfterConsume(viper.GetBool("stop-after-consume")),
	)

	if err := app.Run(cfg); err != nil {
		color.Red("Application error: %v", err)
		os.Exit(1)
	}
}

func validateRequiredFields() error {
	// Validate URL
	urlStr := viper.GetString("url")
	if urlStr == "" {
		return fmt.Errorf("RabbitMQ URL is required")
	}
	if _, err := url.Parse(fmt.Sprintf("%s://%s", getProtocol(), urlStr)); err != nil {
		return fmt.Errorf("invalid RabbitMQ URL: %v", err)
	}

	// Validate queue
	queue := viper.GetString("queue")
	if queue == "" {
		return fmt.Errorf("queue name is required")
	}

	// Validate output file
	output := viper.GetString("output")
	if output == "" {
		return fmt.Errorf("output file name is required")
	}

	// Validate file mode
	fileMode := viper.GetString("file-mode")
	if fileMode != "append" && fileMode != "overwrite" {
		return fmt.Errorf("invalid file mode: must be 'append' or 'overwrite'")
	}

	// Validate virtual host (optional, but if provided, shouldn't start with '/')
	virtualHost := viper.GetString("virtualhost")
	if virtualHost != "" && strings.HasPrefix(virtualHost, "/") {
		return fmt.Errorf("virtual host should not start with '/'")
	}

	return nil
}

func getProtocol() string {
	if viper.GetBool("amqps") {
		return "amqps"
	}
	return "amqp"
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}
