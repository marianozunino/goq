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

var validWriters = []string{"file", "console"}

var rootCmd = &cobra.Command{
	Use:   "goq",
	Short: "A tool to dump RabbitMQ messages to a file",
	Long: logo + `

This application connects to a RabbitMQ server, consumes messages from a specified queue,
and writes them to a file while keeping the messages in the queue.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		run(cmd, args)
		return nil
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Define the flags
	rootCmd.PersistentFlags().StringP("url", "u", "", "RabbitMQ URL (e.g., localhost:5672)")
	rootCmd.PersistentFlags().StringP("exchange", "e", "", "RabbitMQ exchange name")
	rootCmd.PersistentFlags().BoolP("amqps", "s", false, "Use AMQPS instead of AMQP")
	rootCmd.PersistentFlags().StringP("virtualhost", "v", "", "RabbitMQ virtual host")
	rootCmd.PersistentFlags().BoolP("skip-tls-verify", "k", false, "Skip TLS certificate verification (insecure)")

	rootCmd.PersistentFlags().StringP("writer", "w", "file", "Output writer type (console or file)")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output file name (required when writer is 'file')")
	rootCmd.PersistentFlags().StringP("file-mode", "m", "overwrite", "File mode (append or overwrite, only valid for file writer)")

	rootCmd.PersistentFlags().BoolP("pretty-print", "p", false, "Pretty print JSON messages")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", xdg.ConfigHome+"/goq/goq.yaml", "config file")

	rootCmd.PersistentFlags().BoolP("full-message", "f", false, "Print full message")
	rootCmd.PersistentFlags().StringSliceP("include-patterns", "i", []string{}, "Include messages containing these patterns")
	rootCmd.PersistentFlags().StringSliceP("exclude-patterns", "x", []string{}, "Exclude messages containing these patterns")
	rootCmd.PersistentFlags().StringP("json-filter", "j", "", "JSON filter expression")
	rootCmd.PersistentFlags().IntP("max-message-size", "z", -1, "Maximum message size in bytes")
	rootCmd.PersistentFlags().StringP("regex-filter", "R", "", "Regex pattern to filter messages")

	rootCmd.AddGroup(&cobra.Group{
		ID:    "available-commands",
		Title: "Available Commands:",
	})

	// Add flag validation
	rootCmd.PreRunE = validateFlags

	// Allow viper to bind flags
	viper.BindPFlags(rootCmd.PersistentFlags())
}

// validateFlags ensures all flags are correctly set
func validateFlags(cmd *cobra.Command, args []string) error {
	writer, _ := cmd.Flags().GetString("writer")
	output, _ := cmd.Flags().GetString("output")
	fileMode, _ := cmd.Flags().GetString("file-mode")

	// Validate writer
	if !isValidWriter(writer) {
		return fmt.Errorf("invalid writer type '%s', must be one of: %v", writer, validWriters)
	}

	// Validate output file for file writer
	if writer == "file" && output == "" {
		return fmt.Errorf("output file is required when using file writer")
	}

	// Validate file mode
	if fileMode != "append" && fileMode != "overwrite" {
		return fmt.Errorf("invalid file mode '%s': must be 'append' or 'overwrite'", fileMode)
	}

	return nil
}

// isValidWriter checks if the writer is valid
func isValidWriter(writer string) bool {
	for _, valid := range validWriters {
		if writer == valid {
			return true
		}
	}
	return false
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
	switch cmd.Use {
	case "dump", "monitor":
		if err := validateRequiredFields(); err != nil {
			color.Red("Validation error: %v", err)
			os.Exit(1)
		}
	}
}

func validateRequiredFields() error {
	urlStr := viper.GetString("url")
	if urlStr == "" {
		return fmt.Errorf("RabbitMQ URL is required")
	}
	if _, err := url.Parse(fmt.Sprintf("%s://%s", getProtocol(), urlStr)); err != nil {
		return fmt.Errorf("invalid RabbitMQ URL: %v", err)
	}

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
