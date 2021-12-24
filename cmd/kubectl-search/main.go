// Package main provides the CLI entrypoint for the kubectl plugin
// kubectl-search; nothing else.
package main

import (
	"github.com/kubideh/kubesearch/client"
)

func main() {
	client.New().Run()
}
