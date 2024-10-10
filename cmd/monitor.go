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

var routingKeys []string

func NewMonitorCmd() *cobra.Command {
	monitorCmd := &cobra.Command{
		Use:     "monitor",
		Aliases: []string{"mon"},
		GroupID: "available-commands",
		Short:   "Monitor RabbitMQ messages using routing keys and a temporary queue.",
		Long: `Monitor RabbitMQ messages by consuming from a temporary queue that listens to specified routing keys.
This command captures and dumps the received messages to a file for analysis or processing.`,
		Example: `goq monitor -r "key1,key2" -o output.txt -s -v my_vhost`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.New(
				config.WithRabbitMQURL(fmt.Sprintf("%s://%s/%s", getProtocol(), viper.GetString("url"), viper.GetString("virtualhost"))),
				config.WithExchange(viper.GetString("exchange")),
				config.WithOutputFile(viper.GetString("output")),
				config.WithUseAMQPS(viper.GetBool("amqps")),
				config.WithVirtualHost(viper.GetString("virtualhost")),
				config.WithSkipTLSVerification(viper.GetBool("skip-tls-verify")),
				config.WithAutoAck(viper.GetBool("auto-ack")),
				config.WithFileMode(viper.GetString("file-mode")),
				config.WithPrettyPrint(viper.GetBool("pretty-print")),
				config.WithRoutingKeys(routingKeys),
			)

			return app.Monitor(cfg)
		},
	}

	// Allow to pass an array of routing keys to monitor
	monitorCmd.Flags().SortFlags = false
	monitorCmd.Flags().StringSliceVarP(&routingKeys, "routing-keys", "r", nil, "List of routing keys to monitor")
	monitorCmd.MarkFlagRequired("routing-keys")

	return monitorCmd
}

func init() {
	rootCmd.AddCommand(NewMonitorCmd())
}
