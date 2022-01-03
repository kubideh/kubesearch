// Package main provides the CLI entrypoint for the kubectl plugin
// kubectl-search; nothing else.
package main

import (
	"fmt"
	"os"

	"github.com/kubideh/kubesearch/client"
)

func main() {
	aClient := client.ConfigureDefault()

	if err := aClient.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
