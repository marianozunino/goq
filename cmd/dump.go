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

	app "github.com/marianozunino/goq/internal"
	"github.com/marianozunino/goq/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDumpCmd() *cobra.Command {
	var queue string
	var autoAck bool
	var stopAfterConsume bool

	dumpCmd := &cobra.Command{
		Use:   "dump",
		Short: "Dump messages from a RabbitMQ queue to a file.",
		Long: `Dump messages from the specified RabbitMQ queue to a file.
This command provides options for automatically acknowledging messages, controlling when consumption stops, and configuring file output behavior.
Messages can be captured from an AMQP or AMQPS RabbitMQ server, with flexible TLS and virtual host settings.`,
		Example: `goq dump -q my_queue -o output.txt -a -c`,
		GroupID: "available-commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.New(

				config.WithRabbitMQURL(fmt.Sprintf("%s://%s/%s", getProtocol(), viper.GetString("url"), viper.GetString("virtualhost"))),
				config.WithExchange(viper.GetString("exchange")),
				config.WithQueue(queue),
				config.WithOutputFile(viper.GetString("output")),
				config.WithUseAMQPS(viper.GetBool("amqps")),
				config.WithVirtualHost(viper.GetString("virtualhost")),
				config.WithSkipTLSVerification(viper.GetBool("skip-tls-verify")),
				config.WithAutoAck(viper.GetBool("auto-ack")),
				config.WithFileMode(viper.GetString("file-mode")),
				config.WithStopAfterConsume(viper.GetBool("stop-after-consume")),
				config.WithPrettyPrint(viper.GetBool("pretty-print")),
			)

			return app.Dump(cfg)
		},
	}

	dumpCmd.Flags().SortFlags = false
	dumpCmd.Flags().StringVarP(&queue, "queue", "q", "", "RabbitMQ queue name")
	dumpCmd.Flags().BoolVarP(&autoAck, "auto-ack", "a", false, "Auto ack messages")
	dumpCmd.Flags().BoolVarP(&stopAfterConsume, "stop-after-consume", "c", false, "Stop after consuming messages")

	dumpCmd.MarkFlagRequired("queue")

	return dumpCmd
}

func init() {
	rootCmd.AddCommand(NewDumpCmd())
}
