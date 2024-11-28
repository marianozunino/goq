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
	"os"

	"github.com/fatih/color"
	"github.com/marianozunino/goq/pkg/config"
	"github.com/marianozunino/goq/pkg/validation"
	"github.com/spf13/cobra"
)

var (
	cfgFile        string
	validWriters   = []string{"file", "console"}
	validFileModes = []string{"append", "overwrite"}
	logo           = `
  ______    ______    ______
 /\  ___\  /\  __ \  /\  __ \
 \ \ \__ \ \ \ \/\ \ \ \ \/\_\
  \ \_____\ \ \_____\ \ \___\_\
   \/_____/  \/_____/  \/___/_/ ` + VersionFromBuild()
)

var RootCmd = &cobra.Command{
	Use:   "goq",
	Short: "A tool to dump RabbitMQ messages to a file",
	Long:  logo + "\n\nThis application connects to a RabbitMQ server and dumps queue messages to a file.",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validation.ValidateInput(); err != nil {
			color.Red("Validation error: %v", err)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	cobra.OnInitialize(config.InitConfig)
	config.SetupFlags(RootCmd.PersistentFlags(), validWriters, validFileModes)
	RootCmd.AddGroup(&cobra.Group{
		ID:    "available-commands",
		Title: "Available Commands:",
	})
	RootCmd.AddCommand(NewDumpCmd(), NewMonitorCmd(), NewConfigureCommand(), NewUpdateCmd())
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}
