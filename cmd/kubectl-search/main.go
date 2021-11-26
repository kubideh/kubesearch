// Package main provides the CLI entrypoint for the kubectl plugin
// kubectl-search; nothing else.
package main

import "k8s.io/klog/v2"

func main() {
	klog.Infoln("kubectl-search")
}
