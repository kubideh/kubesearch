package app

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

// kubeClientFrom returns Kubernetes client objects
// (`kubernetes.Clientset`) from the configuration given by `flags`.
func kubeClientFrom(flags immutableFlags) *kubernetes.Clientset {
	// use the current context in kubeConfig
	config, err := clientcmd.BuildConfigFromFlags("", flags.KubeConfig())

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
