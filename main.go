// Package main provides the CLI entrypoint for kubesearch; nothing
// else.
package main

import (
	"flag"
	"net/http"
	"path/filepath"

	"github.com/kubideh/kubesearch/search"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		klog.Fatalln(err)
	}

	// create the clientset
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalln(err)
	}

	// create the Controller to be used by the search API handler
	controller := search.NewController(client)
	search.SetControllerRef(controller)
	cancel := controller.Start()
	defer cancel()

	klog.Infoln("Listening on :8080")
	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		klog.Fatalln(err)
	}
}
