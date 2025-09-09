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

// NewDumpCmd creates the `dump` command.
func NewDumpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump",
		Short: "Dump messages from a RabbitMQ queue",
		Long:  "Dump messages from a specified RabbitMQ queue with flexible filtering and output options.",
		Example: `  # Dump messages from a queue to file
  goq dump -q "my_queue" -o messages.json -p

  # Dump with filtering and auto-acknowledge
  goq dump -q "orders" -a -i "urgent" -o urgent_orders.log

  # Dump from secure connection with full message details
  goq dump -q "events" -s -k -f -o events_full.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Dump(config.CreateCommonConfig(cmd))
		},
	}

	cmd.Flags().StringP("queue", "q", "", "RabbitMQ queue name (required)")
	cmd.Flags().BoolP("auto-ack", "a", false, "Automatically acknowledge messages")
	cmd.Flags().BoolP("stop-after-consume", "c", false, "Stop after consuming messages")
	cmd.Flags().BoolP("full-message", "f", false, "Print complete message details")
	cmd.MarkFlagRequired("queue")

	return cmd
}
