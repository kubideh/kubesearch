// Package index provides an inverted index that supports fulltext
// search.
package index

import (
	"fmt"
	"sort"
	"sync"
)

// Posting represents an object Key and the kind of object the ID
// references.
type Posting struct {
	Key  string
	Kind string
}

// DocID is the document identifier, and it's a string with the
// form <Kind>/<Optional namespace>/<Object name>.
func (p Posting) DocID() string {
	return fmt.Sprintf("%s/%s", p.Kind, p.Key)
}

type PostingsList []Posting

func (p PostingsList) Len() int {
	return len(p)
}

func (p PostingsList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PostingsList) Less(i, j int) bool {
	return p[i].DocID() < p[j].DocID()
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
		postings := idx.index[t]

		if contains(postings, posting) {
			return
		}

		postings = append(postings, posting)
		sort.Sort(PostingsList(postings))
		idx.index[t] = postings
	}
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
