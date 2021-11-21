package search

import (
	"k8s.io/client-go/tools/cache"
)

var singletonImformer cache.SharedIndexInformer

// Informer returns the global K8s informer.
func Informer() cache.SharedIndexInformer {
	return singletonImformer
}

// SetInformer replaces the global K8s informer.
func SetInformer(informer cache.SharedIndexInformer) {
	singletonImformer = informer
}
