// Package main provides the CLI entrypoint for the kubectl plugin
// kubectl-search; nothing else.
package main

import (
	"fmt"
	"os"

	"github.com/kubideh/kubesearch/search/api"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "You must specify a query.")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Use \"kubectl search <query>\".")
		os.Exit(1)
	}

	if err := api.Search("http://localhost:8080", os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
