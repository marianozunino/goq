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
	app "github.com/marianozunino/goq/internal"
	"github.com/marianozunino/goq/pkg/config"
	"github.com/spf13/cobra"
)

// NewMonitorCmd creates the `monitor` command.
func NewMonitorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "monitor",
		Aliases: []string{"mon"},
		Short:   "Monitor RabbitMQ messages using routing keys",
		Long:    "Monitor RabbitMQ messages by consuming from a temporary queue with specified routing keys.",
		Example: `  # Monitor all messages from an exchange
  goq monitor -K "#" -e "my_exchange" -w console -p

  # Monitor specific routing keys with filtering
  goq monitor -K "user.created,user.updated" -e "events" -i "admin" -o users.log

  # Monitor with secure connection
  goq monitor -K "order.*" -e "orders" -s -k -u "rabbitmq.example.com:5671"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Monitor(config.CreateCommonConfig(cmd))
		},
	}

	cmd.Flags().StringSliceP("routing-keys", "K", nil, "List of routing keys to monitor (required)")
	cmd.Flags().BoolP("auto-ack", "a", false, "Automatically acknowledge messages")
	cmd.MarkFlagRequired("routing-keys")

	return cmd
}
