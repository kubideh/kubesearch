package finder

import (
	"fmt"

	"k8s.io/client-go/tools/cache"
)

// XXX This finder is actually a Gateway. Maybe refactor into Active Records.s

// Object is a (Key, Item) pair.
type Object struct {
	Key  Key
	Item interface{}
}

type Key struct {
	StoredObjectKey string
	K8sResourceKind string
}

type FindAllFunc func(keys []Key) ([]Object, error)

// Create returns the default functor that finds all objects for
// the given Kubernetes object store `store`.
func Create(store map[string]cache.Store) FindAllFunc {
	return func(keys []Key) ([]Object, error) {
		var results []Object

		for _, k := range keys {
			item, exists, err := findOne(store, k.K8sResourceKind, k.StoredObjectKey)

			if err != nil {
				return results, err
			}

			if !exists {
				return results, fmt.Errorf("missing object for key %v", k.StoredObjectKey)
			}

			results = append(results, Object{Key: k, Item: item})
		}

		return results, nil
	}
}

func findOne(store map[string]cache.Store, kind, key string) (item interface{}, exists bool, err error) {
	return store[kind].GetByKey(key)
}
