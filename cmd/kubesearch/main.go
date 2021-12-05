// Package main provides the CLI entrypoint for kubesearch; nothing
// else.
package main

import (
	"github.com/kubideh/kubesearch/search/app"
)

func main() {
	app.New().Run()
}
