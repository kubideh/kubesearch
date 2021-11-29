package search

import (
	"bufio"
	"context"
	"strings"
	"unicode/utf8"

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
	for kind, informer := range c.informers {
		go indexObjects(informer.queue, index, kind)
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
			index.Put(dnsSubdomainNamesTokenizer(namespace(key)), Posting{Key: keyString(key), Kind: kind})
		}

		index.Put(dnsSubdomainNamesTokenizer(name(key)), Posting{Key: keyString(key), Kind: kind})

		// XXX Support indexing annotations and labels

		key, shutdown = queue.Get()
	}

	klog.Infof("Shutting down %s queue", kind)
}

// The DNS Subdomain Name tokenizer follows the rules for naming
// objects in Kubernetes (https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names).
// In addition to tokenizing on hyphens or dots, the exact name
// is also returned as the first token. For example, for the name
// `dns.sub-domain.name`, the following tokens are returned:
// `dns`, `sub`, `domain`, `name`, and `dns.sub-domain.name`.
func dnsSubdomainNamesTokenizer(text string) (results []string) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(scan)

	for scanner.Scan() {
		results = append(results, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		klog.Warningln("scanner error: ", err)
	}

	if len(results) > 1 {
		results = append(results, text)
	}

	return
}

// scan is a split function for a Scanner that returns UTF-8 tokens
// split on dots or hyphens. This algorithm is taken from bufio.ScanWords.
func scan(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if r != '.' && r != '-' {
			break
		}
	}

	// Scan until dot or hyphen, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if r == '.' || r == '-' {
			return i + width, data[start:i], nil
		}
	}

	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	// Request more data.
	return start, nil, nil
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
