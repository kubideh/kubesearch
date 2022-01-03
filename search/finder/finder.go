package finder

import (
	"fmt"

	"k8s.io/client-go/tools/cache"
)

// XXX This finder is actually a Gateway. Maybe refactor into Active Records.s

// K8sObject is a (Key, Item) pair.
type K8sObject struct {
	Key  Key
	Item interface{}
}

type Key struct {
	StoredObjectKey string
	K8sResourceKind string
}

type FindAllFunc func(keys []Key) ([]K8sObject, error)

// Create returns the default functor that finds all objects for
// the given Kubernetes object store `store`.
func Create(store map[string]cache.Store) FindAllFunc {
	return func(keys []Key) ([]K8sObject, error) {
		var results []K8sObject

		for _, k := range keys {
			item, exists, err := findOne(store, k.K8sResourceKind, k.StoredObjectKey)

			if err != nil {
				return results, err
			}

			if !exists {
				// XXX: Ignore because the object might have been deleted.
				return results, fmt.Errorf("missing object for key %v", k.StoredObjectKey)
			}

			object := K8sObject{
				Key:  k,
				Item: item,
			}

			results = append(results, object)
		}

		return results, nil
	}
}

func findOne(store map[string]cache.Store, kind, key string) (item interface{}, exists bool, err error) {
	return store[kind].GetByKey(key)
}
