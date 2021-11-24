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
	"k8s.io/client-go/util/workqueue"
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

	// create a workqueue and index to be used by the seach API handler
	index := search.NewIndex()
	podQueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "pod-queue")
	go search.IndexObjects(podQueue, index)

	// create the informer to be used by the search API handler
	informerFactory := informers.NewSharedInformerFactory(clientset, 0)
	podInformer := search.NewPodInformer(informerFactory, podQueue)
	controller := search.NewController(podInformer, podQueue, index)
	search.SetControllerRef(controller)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	informerFactory.Start(ctx.Done())

	klog.Info("Listening on :8080")
	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		klog.Fatalln(err)
	}
}
