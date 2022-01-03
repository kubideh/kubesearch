// Package main provides the CLI entrypoint for kubesearch; nothing
// else.
package main

import (
	"github.com/kubideh/kubesearch/cmd/kubesearch/app"
	"k8s.io/klog/v2"
)

func main() {
	anApp := app.ConfigureDefault()

	if err := anApp.Run(); err != nil {
		klog.Fatalln(err)
	}
}
