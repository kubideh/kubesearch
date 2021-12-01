// Package controller provides a Kubernetes controller that watches
// for changes to Kubernetes objects and updates the an inverted
// index with those changes.
package controller

import (
	"context"
	"strings"

	"github.com/kubideh/kubesearch/search/index"
	"github.com/kubideh/kubesearch/search/tokenizer"
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

// Informer binds an informer and a workqueue.
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
func (c *Controller) Start(idx *index.Index) context.CancelFunc {
	for kind, informer := range c.informers {
		go indexObjects(informer.queue, idx, kind)
	}

	ctx, cancel := context.WithCancel(context.Background())

	c.informerFactory.Start(ctx.Done())

	c.informerFactory.WaitForCacheSync(ctx.Done())

	return cancel
}

// New returns Controller objects.
func New(client kubernetes.Interface) *Controller {
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

func indexObjects(queue workqueue.RateLimitingInterface, idx *index.Index, kind string) {
	key, shutdown := queue.Get()

	for !shutdown {
		if namespace(key) != "" {
			idx.Put(tokenizer.Tokenize(namespace(key)), index.Posting{Key: keyString(key), Kind: kind})
		}

		idx.Put(tokenizer.Tokenize(name(key)), index.Posting{Key: keyString(key), Kind: kind})

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
