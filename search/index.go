package search

import "sync"

// Posting represents an object Key and the kind of object the ID
// references.
type Posting struct {
	Key  string
	Kind string
}

// InvertedIndex maps terms to object keys.
type InvertedIndex struct {
	index map[string][]Posting
	mutex sync.RWMutex
}

// Put adds a Posting to the search index for the given term.
func (idx *InvertedIndex) Put(term string, doc Posting) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	idx.index[term] = append(idx.index[term], doc)
}

// Get looks up Postings in the search index using the given term.
func (idx *InvertedIndex) Get(term string) (result []Posting, found bool) {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()

	result, found = idx.index[term]

	return
}

// NewIndex returns InvertedIndex objects.
func NewIndex() *InvertedIndex {
	return &InvertedIndex{
		index: make(map[string][]Posting),
	}
}
