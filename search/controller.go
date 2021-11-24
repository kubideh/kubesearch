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
	informer        cache.SharedIndexInformer
	queue           workqueue.RateLimitingInterface
	index           InvertedIndex
}

// GetPodByKey returns the pod for the given key.
func (c *Controller) GetPodByKey(key string) (item interface{}, exists bool, err error) {
	return c.informer.GetStore().GetByKey(key)
}

// Index returns a reference to the inverted search index.
func (c *Controller) Index() InvertedIndex {
	return c.index
}

// Start this controller. The caller should defer the call to the
// return cancel function.
func (c *Controller) Start() context.CancelFunc {
	go indexObjects(c.queue, c.index)
	ctx, cancel := context.WithCancel(context.Background())
	c.informerFactory.Start(ctx.Done())
	return cancel
}

// NewController returns Controller objects.
func NewController(client kubernetes.Interface) *Controller {
	factory := informers.NewSharedInformerFactory(client, 0)
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "pod-queue")
	informer := newPodInformer(factory, queue)
	index := NewIndex()

	return &Controller{
		informerFactory: factory,
		informer:        informer,
		queue:           queue,
		index:           index,
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

func indexObjects(queue workqueue.RateLimitingInterface, index InvertedIndex) {
	key, shutdown := queue.Get()
	for !shutdown {
		var name string

		metadata := strings.Split(key.(string), "/")
		if len(metadata) == 1 {
			name = metadata[0]
		} else {
			name = metadata[1]
		}

		index[name] = key.(string)

		key, shutdown = queue.Get()
	}
	klog.Infoln("Shutting down pod-queue")
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
