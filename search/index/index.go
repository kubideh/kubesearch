// Package index provides an inverted index that supports fulltext
// search.
package index

import (
	"sort"
	"sync"
)

// Index maps terms to object keys.
type Index struct {
	index map[string][]Posting
	mutex sync.RWMutex
}

// Put adds a posting to the search index for each of the given
// terms. The frequency of terms in the posting is incremented.
func (idx *Index) Put(terms []string, posting Posting) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	for _, t := range terms {
		idx.putOne(t, posting)
	}
}

func (idx *Index) putOne(term string, posting Posting) {
	postings := idx.index[term]

	if contains(postings, posting) {
		return
	}

	posting.Frequency = posting.TermFrequency(term)
	postings = append(postings, posting)

	sort.Sort(PostingsList(postings))

	idx.index[term] = postings
}

func contains(postings []Posting, item Posting) bool {
	for _, p := range postings {
		if p.DocID() == item.DocID() {
			return true
		}
	}

	return false
}

// Get looks up a posting list in the search index using the given term.
func (idx *Index) Get(term string) []Posting {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()

	return idx.index[term]
}

// New returns InvertedIndex objects.
func New() *Index {
	return &Index{
		index: make(map[string][]Posting),
	}
}
