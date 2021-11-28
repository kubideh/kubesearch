package search

import "sync"

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

// Put adds a posting to the search index for the given term.
func (idx *Index) Put(term string, posting Posting) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	idx.index[term] = append(idx.index[term], posting)
}

// Get looks up a posting list in the search index using the given term.
func (idx *Index) Get(term string) (result []Posting, found bool) {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()

	result, found = idx.index[term]

	return
}

// NewIndex returns InvertedIndex objects.
func NewIndex() *Index {
	return &Index{
		index: make(map[string][]Posting),
	}
}

// DoIndex indexes the given text for the given posting.
func DoIndex(index *Index, text string, posting Posting) {
	index.Put(text, posting)
}
