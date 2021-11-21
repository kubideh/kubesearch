// Package main provides the CLI entrypoint for kubesearch; nothing
// else.
package main

import (
	"context"
	"flag"
	"net/http"
	"path/filepath"

	"github.com/kubideh/kubesearch/search"
	"k8s.io/client-go/informers"
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
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalln(err)
	}

	// create the informer to be used by the search API handler
	informerFactory := informers.NewSharedInformerFactory(clientset, 0)
	search.SetInformer(informerFactory.Core().V1().Pods().Informer())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	informerFactory.Start(ctx.Done())

	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		klog.Fatalln(err)
	}
}
