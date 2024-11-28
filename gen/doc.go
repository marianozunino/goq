package main

import (
	"github.com/marianozunino/goq/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	app := cmd.RootCmd
	app.DisableAutoGenTag = true
	doc.GenMarkdownTree(cmd.RootCmd, "./docs")
}
