// Package main provides the CLI entrypoint for kubesearch; nothing
// else.
package main

import (
	"github.com/kubideh/kubesearch/app"
)

func main() {
	app.Create().Run()
}
