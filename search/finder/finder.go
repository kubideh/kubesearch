package finder

import (
	"fmt"

	"github.com/kubideh/kubesearch/search/index"
	"k8s.io/client-go/tools/cache"
)

// XXX This finder is actually a Gateway. Maybe refactor into Active Records.s

// Object is a (Posting, Item) pair.
type Object struct {
	Posting index.Posting
	Item    interface{}
}

type FindAllFunc func(postings []index.Posting) ([]Object, error)

// Create returns the default functor that finds all objects for
// the given Kubernetes object store `store`.
func Create(store map[string]cache.Store) FindAllFunc {
	return func(postings []index.Posting) ([]Object, error) {
		var results []Object

		for _, p := range postings {
			item, exists, err := findOne(store, p.K8sResourceKind, p.StoredObjectKey)

			if err != nil {
				return results, err
			}

			if !exists {
				return results, fmt.Errorf("missing object for key %v", p.StoredObjectKey)
			}

			results = append(results, Object{Posting: p, Item: item})
		}

		return results, nil
	}
}

func findOne(store map[string]cache.Store, kind, key string) (item interface{}, exists bool, err error) {
	return store[kind].GetByKey(key)
}
