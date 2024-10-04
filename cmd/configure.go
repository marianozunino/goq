/*
Copyright © 2024 Mariano Zunino <marianoz@posteo.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

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
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Create a sample configuration file",
	Long:  `Create a sample configuration file for rabbitmq-dumper in the default location.`,
	Run:   runConfigure,
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

func runConfigure(cmd *cobra.Command, args []string) {
	exampleConfig := `# Configuration for rabbitmq-dumper

# RabbitMQ server URL
url: "localhost:5672"

# RabbitMQ exchange name
exchange: "my_exchange"

# RabbitMQ queue name
queue: "my_queue"

# Output file name
output: "messages.txt"

# Use AMQPS instead of AMQP
amqps: false

# RabbitMQ virtual host
virtualhost: ""

# Skip TLS certificate verification (insecure)
skip-tls-verify: false

# Automatically acknowledge messages
auto-ack: false

# File mode (append or overwrite)
file-mode: "overwrite"

# Stop consuming after getting all messages from the queue
stop-after-consume: false
`

	configPath, err := xdg.ConfigFile("goq/goq.yaml")
	if err != nil {
		log.Fatalf("Failed to determine config file path: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}

	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Configuration file already exists at %s\n", configPath)
		fmt.Print("Do you want to overwrite it? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Configuration generation cancelled.")
			return
		}
	}

	err = os.WriteFile(configPath, []byte(exampleConfig), 0o644)
	if err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}

	configPath, _ = filepath.Abs(configPath)

	fmt.Printf("Configuration file generated: %s\n", configPath)
	fmt.Println("You can now edit this file to customize your settings.")
}
