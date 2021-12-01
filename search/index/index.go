// Package index provides an inverted index that supports fulltext
// search.
package index

import (
	"sync"
)

// Posting represents an object Key and the kind of object the ID
// references.
type Posting struct {
	Key  string
	Kind string
}

// Index maps terms to object keys.
type Index struct {
	index map[string][]Posting
	mutex sync.RWMutex
}

// Put adds a posting to the search index for each of the given
// terms.
func (idx *Index) Put(terms []string, posting Posting) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	for _, t := range terms {
		idx.index[t] = append(idx.index[t], posting)
	}
}

// Get looks up a posting list in the search index using the given term.
func (idx *Index) Get(term string) (result []Posting, found bool) {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()

	result, found = idx.index[term]

	return
}

// New returns InvertedIndex objects.
func New() *Index {
	return &Index{
		index: make(map[string][]Posting),
	}
}
