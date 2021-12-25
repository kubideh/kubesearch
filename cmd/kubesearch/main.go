// Package main provides the CLI entrypoint for kubesearch; nothing
// else.
package main

import (
	"github.com/kubideh/kubesearch/app"
	"k8s.io/klog/v2"
)

func main() {
	anApp := app.Create()

	if err := anApp.Run(); err != nil {
		klog.Fatalln(err)
	}
}
