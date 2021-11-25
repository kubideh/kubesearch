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
	kubeconfig := defineKubeconfigFlag()
	flag.Parse()

	client := createKubernetesClient(kubeconfig)

	// create the Controller to be used by the search API handler
	controller := search.NewController(client)

	index := search.NewIndex()

	cancel := controller.Start(index)
	defer cancel()

	mux := http.NewServeMux()
	search.RegisterHandler(mux, index, controller.Store())

	klog.Infoln("Listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		klog.Fatalln(err)
	}
}

func defineKubeconfigFlag() (kubeconfig *string) {
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	return
}

func createKubernetesClient(kubeconfig *string) *kubernetes.Clientset {
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

	return client
}
