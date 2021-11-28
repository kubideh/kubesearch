package search

import (
	"strings"
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

// IndexDNSSubdomainNames indexes the given text for the given
// posting using the rules given by https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names.
func IndexDNSSubdomainNames(index *Index, text string, posting Posting) {
	if text == "" {
		return
	}

	if len(text) > 253 {
		text = text[0:253]
	}

	// XXX replace with Scanner or regex that tokenizes into words
	// including words separated by dots or hyphens. the set of
	// tokens should be a powerset.

	tokens := make(map[string]struct{})

	tokens[text] = struct{}{} // The entire text is a token.

	for _, t := range strings.Split(text, ".") {
		tokens[t] = struct{}{}

		for _, s := range strings.Split(t, "-") {
			tokens[s] = struct{}{}
		}
	}

	for _, t := range strings.Split(text, "-") {
		tokens[t] = struct{}{}

		for _, s := range strings.Split(t, ".") {
			tokens[s] = struct{}{}
		}
	}

	for t := range tokens {
		index.Put(t, posting)
	}
}
