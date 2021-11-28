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
	informerFactory informers.SharedInformerFactory
	informers       map[string]Informer
}

type Informer struct {
	informer cache.SharedIndexInformer
	queue    workqueue.RateLimitingInterface
}

// Store returns the object store.
func (c *Controller) Store() map[string]cache.Store {
	result := make(map[string]cache.Store)

	for k := range c.informers {
		result[k] = c.informers[k].informer.GetStore()
	}

	return result
}

// Start this controller. The caller should defer the call to the
// return cancel function.
func (c *Controller) Start(index *Index) context.CancelFunc {
	for k := range c.informers {
		go indexObjects(c.informers[k].queue, index, k)
	}

	ctx, cancel := context.WithCancel(context.Background())

	c.informerFactory.Start(ctx.Done())

	c.informerFactory.WaitForCacheSync(ctx.Done())

	return cancel
}

// NewController returns Controller objects.
func NewController(client kubernetes.Interface) *Controller {
	factory := informers.NewSharedInformerFactory(client, 0)

	// XXX Support the creation of informers by the caller.

	return &Controller{
		informerFactory: factory,
		informers: map[string]Informer{
			"Deployment": newInformer(factory.Apps().V1().Deployments().Informer(), "Deployment-queue"),
			"Pod":        newInformer(factory.Core().V1().Pods().Informer(), "Pod-queue"),
		},
	}
}

func indexObjects(queue workqueue.RateLimitingInterface, index *Index, kind string) {
	key, shutdown := queue.Get()

	for !shutdown {
		if namespace(key) != "" {
			// XXX The name of your namespace must be a valid DNS label.
			// https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-label-names
			IndexDNSSubdomainNames(index, namespace(key), Posting{Key: keyString(key), Kind: kind})
		}

		IndexDNSSubdomainNames(index, name(key), Posting{Key: keyString(key), Kind: kind})

		// XXX Support indexing annotations and labels

		key, shutdown = queue.Get()
	}

	klog.Infof("Shutting down %s queue", kind)
}

func keyString(key interface{}) string {
	return key.(string)
}

func namespace(key interface{}) (result string) {
	metadata := strings.Split(keyString(key), "/")

	if len(metadata) > 1 {
		result = metadata[0]
	}

	return
}

func name(key interface{}) (result string) {
	metadata := strings.Split(keyString(key), "/")

	if len(metadata) == 1 {
		result = metadata[0]
	} else {
		result = metadata[1]
	}

	return
}

func newInformer(informer cache.SharedIndexInformer, name string) Informer {
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), name)

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

	return Informer{
		informer: informer,
		queue:    queue,
	}
}
