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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		run(cmd, args)
		return nil
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("url", "u", "", "RabbitMQ URL (e.g., localhost:5672)")
	rootCmd.PersistentFlags().StringP("exchange", "e", "", "RabbitMQ exchange name")
	rootCmd.PersistentFlags().StringP("output", "o", "messages.txt", "Output file name")
	rootCmd.PersistentFlags().BoolP("amqps", "s", false, "Use AMQPS instead of AMQP")
	rootCmd.PersistentFlags().StringP("virtualhost", "v", "", "RabbitMQ virtual host")
	rootCmd.PersistentFlags().BoolP("skip-tls-verify", "k", false, "Skip TLS certificate verification (insecure)")
	rootCmd.PersistentFlags().StringP("file-mode", "m", "overwrite", "File mode (append or overwrite)")
	rootCmd.PersistentFlags().BoolP("pretty-print", "p", false, "Pretty print JSON messages")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", xdg.ConfigHome+"/goq/goq.yaml", "config file")

	rootCmd.AddGroup(&cobra.Group{
		ID:    "available-commands",
		Title: "Available Commands:",
	})

	// allow to attatch to multiple
	viper.BindPFlags(rootCmd.PersistentFlags())
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
	viper.ReadInConfig()
}

func run(cmd *cobra.Command, args []string) {
	// Validate required fields
	switch cmd.Use {
	case "dump", "monitor":
		if err := validateRequiredFields(); err != nil {
			color.Red("Validation error: %v", err)
			os.Exit(1)
		}
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
