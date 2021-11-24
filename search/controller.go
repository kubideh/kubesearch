package search

import (
	"strings"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

// Controller is an informer, a workqueue, and an inverted index.
type Controller struct {
	informer cache.SharedIndexInformer
	queue    workqueue.RateLimitingInterface
	index    InvertedIndex
}

// NewController returns Controller objects.
func NewController(podInformer cache.SharedIndexInformer, podQueue workqueue.RateLimitingInterface, index InvertedIndex) *Controller {
	return &Controller{
		informer: podInformer,
		queue:    podQueue,
		index:    index,
	}
}

var singletonController *Controller

// ControllerRef returns the global K8s Controller.
func ControllerRef() *Controller {
	return singletonController
}

// SetControllerRef replaces the global K8s Controller.
func SetControllerRef(controller *Controller) {
	singletonController = controller
}

// IndexObjects indexes the objects taken from the given queue.
func IndexObjects(podQueue workqueue.RateLimitingInterface, index InvertedIndex) {
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
}

// NewPodInformer returns informers for Pods.
func NewPodInformer(informerFactory informers.SharedInformerFactory, podQueue workqueue.RateLimitingInterface) cache.SharedIndexInformer {
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
	return podInformer
}
