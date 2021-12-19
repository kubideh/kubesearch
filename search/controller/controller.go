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
	index           *index.Index
	informerFactory informers.SharedInformerFactory
	informers       map[string]informerWorkqueuePair
	tokenizer       tokenizer.TokenizeFunc
}

// New returns Controller objects.
func New(client kubernetes.Interface) *Controller {
	factory := informers.NewSharedInformerFactory(client, 0)

	// XXX Support the creation of informers by the caller.

	return &Controller{
		index:           index.New(),
		informerFactory: factory,
		informers: map[string]informerWorkqueuePair{
			"Deployment": bindInformerToNewWorkqueue(factory.Apps().V1().Deployments().Informer(), "Deployment-queue"),
			"Pod":        bindInformerToNewWorkqueue(factory.Core().V1().Pods().Informer(), "Pod-queue"),
		},
		tokenizer: tokenizer.Tokenizer(),
	}
}

// Index returns the index bound to this Controller.
func (c *Controller) Index() *index.Index {
	return c.index
}

// Store returns the object store.
func (c *Controller) Store() map[string]cache.Store {
	result := make(map[string]cache.Store)

	for k := range c.informers {
		result[k] = c.store(k)
	}

	return result
}

func (c *Controller) store(kind string) cache.Store {
	return c.informers[kind].informer.GetStore()
}

// Start this controller. The caller should defer the call to the
// return cancel function.
func (c *Controller) Start() context.CancelFunc {
	c.startIndexers()

	ctx, cancel := context.WithCancel(context.Background())
	c.informerFactory.Start(ctx.Done())
	c.informerFactory.WaitForCacheSync(ctx.Done())

	return cancel
}

func (c *Controller) startIndexers() {
	for kind, informer := range c.informers {
		startIndexer(informer.queue, c.index, c.tokenizer, kind)
	}
}

func startIndexer(queue workqueue.RateLimitingInterface, idx *index.Index, tokenize tokenizer.TokenizeFunc, kind string) {
	go indexObjects(queue, idx, tokenize, kind)
}

func indexObjects(queue workqueue.RateLimitingInterface, idx *index.Index, tokenize tokenizer.TokenizeFunc, kind string) {
	key, shutdown := queue.Get()

	for !shutdown {
		if namespace(key) != "" {
			idx.Put(tokenize(namespace(key)), index.Posting{Key: keyString(key), Kind: kind})
		}

		idx.Put(tokenize(name(key)), index.Posting{Key: keyString(key), Kind: kind})

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

func bindInformerToNewWorkqueue(informer cache.SharedIndexInformer, name string) informerWorkqueuePair {
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), name)

	addEventHandlerToInformerUsingQueue(informer, queue)

	return newInformerWorkqueuePair(informer, queue)
}

// informerWorkqueuePair binds an informer and a workqueue.
type informerWorkqueuePair struct {
	informer cache.SharedIndexInformer
	queue    workqueue.RateLimitingInterface
}

func newInformerWorkqueuePair(informer cache.SharedIndexInformer, queue workqueue.RateLimitingInterface) informerWorkqueuePair {
	return informerWorkqueuePair{
		informer: informer,
		queue:    queue,
	}
}

func addEventHandlerToInformerUsingQueue(informer cache.SharedIndexInformer, queue workqueue.RateLimitingInterface) {
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
}
