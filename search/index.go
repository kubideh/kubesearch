package search

import "sync"

// InvertedIndex maps terms to object keys.
type InvertedIndex struct {
	index map[string]string
	mutex sync.RWMutex
}

func (idx *InvertedIndex) Put(term, docID string) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()
	idx.index[term] = docID
}

func (idx *InvertedIndex) Get(term string) (string, bool) {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()
	result, found := idx.index[term]
	return result, found
}

// NewIndex returns a Index objects.
func NewIndex() *InvertedIndex {
	return &InvertedIndex{
		index: make(map[string]string),
	}
}
