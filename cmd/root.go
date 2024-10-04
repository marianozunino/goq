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
	"os"

	"github.com/adrg/xdg"
	"github.com/fatih/color"
	"github.com/marianozunino/goq/internal/config"
	"github.com/marianozunino/goq/internal/filewriter"
	"github.com/marianozunino/goq/internal/rmq"
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

	configJSON, err := cfg.PrettyPrint()
	if err != nil {
		log.Fatalf("Failed to print config: %v", err)
	}
	color.Green("Configuration used:")
	fmt.Println(configJSON)

	consumer, err := rmq.NewConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	writer, err := filewriter.NewWriter(cfg)
	if err != nil {
		log.Fatalf("Failed to create file writer: %v", err)
	}
	defer writer.Close()

	msgCount, err := consumer.GetQueueInfo()
	if err != nil {
		log.Fatalf("Failed to get queue info: %v", err)
	}
	color.Cyan("Number of messages in the queue: %d", msgCount)

	msgs, err := consumer.ConsumeMessages()
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	if cfg.StopAfterConsume {
		color.Yellow("Stopping after consuming all messages")
	}

	log.Println("Waiting for messages. To exit press CTRL+C")

	consumedCount := 0
	blue := color.New(color.FgBlue)
	for msg := range msgs {
		err := writer.WriteMessage(string(msg.Body))
		if err != nil {
			log.Printf("Failed to write message: %v", err)
		}

		consumedCount++
		// color.Blue("\rMessages dumped: %d/%d", consumedCount, msgCount)
		blue.Printf("\rMessages dumped: %d/%d", consumedCount, msgCount)
		os.Stdout.Sync() // Flush the output

		if cfg.StopAfterConsume && consumedCount >= msgCount {
			break
		}
	}
	fmt.Println()

	color.Green("All messages have been consumed. Exiting.")
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
