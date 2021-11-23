// Package main provides the CLI entrypoint for kubesearch; nothing
// else.
package main

import (
	"context"
	"flag"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/kubideh/kubesearch/search"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
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
	index := make(map[string]string)
	search.SetIndex(index)

	podQueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "pod-queue")
	go func(podQueue workqueue.RateLimitingInterface) {
		key, shutdown := podQueue.Get()
		for !shutdown {
			var name string

			metadata := strings.Split(key.(string), "/")
			if len(metadata) == 1 {
				name = metadata[0]
			} else {
				name = metadata[1]
			}

			index[name] = key.(string)

			key, shutdown = podQueue.Get()
		}
		klog.Info("Shutting down pod-queue")
	}(podQueue)

	// create the informer to be used by the search API handler
	informerFactory := informers.NewSharedInformerFactory(clientset, 0)

	podInformer := informerFactory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				klog.Error(err)
			} else {
				podQueue.Add(key)
			}
		},
	})
	search.SetInformer(podInformer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	informerFactory.Start(ctx.Done())

	klog.Info("Listening on :8080")
	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		klog.Fatalln(err)
	}
}
