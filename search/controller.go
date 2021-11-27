package search

import (
	"context"
	"strings"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

// Controller is an informer, a workqueue, and an inverted index.
type Controller struct {
	informerFactory    informers.SharedInformerFactory
	deploymentInformer cache.SharedIndexInformer
	deploymentQueue    workqueue.RateLimitingInterface
	podInformer        cache.SharedIndexInformer
	podQueue           workqueue.RateLimitingInterface
}

// Store returns the object store.
func (c *Controller) Store() map[string]cache.Store {
	return map[string]cache.Store{
		"Deployment": c.deploymentInformer.GetStore(),
		"Pod":        c.podInformer.GetStore(),
	}
}

// Start this controller. The caller should defer the call to the
// return cancel function.
func (c *Controller) Start(index *InvertedIndex) context.CancelFunc {
	go indexDeployments(c.deploymentQueue, index)
	go indexPods(c.podQueue, index)
	ctx, cancel := context.WithCancel(context.Background())
	c.informerFactory.Start(ctx.Done())
	c.informerFactory.WaitForCacheSync(ctx.Done())
	return cancel
}

// NewController returns Controller objects.
func NewController(client kubernetes.Interface) *Controller {
	factory := informers.NewSharedInformerFactory(client, 0)

	deploymentQueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "deployment-queue")
	deploymentInformer := newDeploymentInformer(factory, deploymentQueue)

	podQueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "pod-queue")
	podInformer := newPodInformer(factory, podQueue)

	return &Controller{
		informerFactory:    factory,
		deploymentInformer: deploymentInformer,
		deploymentQueue:    deploymentQueue,
		podInformer:        podInformer,
		podQueue:           podQueue,
	}
}

func indexDeployments(queue workqueue.RateLimitingInterface, index *InvertedIndex) {
	key, shutdown := queue.Get()
	for !shutdown {
		var namespace, name string

		metadata := strings.Split(key.(string), "/")
		if len(metadata) == 1 {
			name = metadata[0]
		} else {
			namespace = metadata[0]
			name = metadata[1]
		}

		if namespace != "" {
			index.Put(namespace, Posting{Key: key.(string), Kind: "Deployment"})
		}

		index.Put(name, Posting{Key: key.(string), Kind: "Deployment"})

		key, shutdown = queue.Get()
	}
	klog.Infoln("Shutting down deployment-queue")
}

func indexPods(queue workqueue.RateLimitingInterface, index *InvertedIndex) {
	key, shutdown := queue.Get()
	for !shutdown {
		var namespace, name string

		metadata := strings.Split(key.(string), "/")
		if len(metadata) == 1 {
			name = metadata[0]
		} else {
			namespace = metadata[0]
			name = metadata[1]
		}

		if namespace != "" {
			index.Put(namespace, Posting{Key: key.(string), Kind: "Pod"})
		}

		index.Put(name, Posting{Key: key.(string), Kind: "Pod"})

		key, shutdown = queue.Get()
	}
	klog.Infoln("Shutting down pod-queue")
}

func newDeploymentInformer(informerFactory informers.SharedInformerFactory, queue workqueue.RateLimitingInterface) cache.SharedIndexInformer {
	informer := informerFactory.Apps().V1().Deployments().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				klog.Errorln(err)
			} else {
				queue.Add(key)
			}
		},
	})
	return informer
}

func newPodInformer(informerFactory informers.SharedInformerFactory, queue workqueue.RateLimitingInterface) cache.SharedIndexInformer {
	informer := informerFactory.Core().V1().Pods().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				klog.Errorln(err)
			} else {
				queue.Add(key)
			}
		},
	})
	return informer
}
