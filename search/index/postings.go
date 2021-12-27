package index

import (
	"fmt"
	"strings"
)

// Posting represents an object Key, the kind of object the ID
// references, and a term Frequency.
type Posting struct {
	Key       string
	Kind      string
	Frequency int
}

// DocID is the document identifier, and it's a string with the
// form <Kind>/<Optional namespace>/<Object name>.
func (p Posting) DocID() string {
	return fmt.Sprintf("%s/%s", p.Kind, p.Key)
}

// TermFrequency returns the number of times term appears in the
// given Posting.
func (p Posting) TermFrequency(term string) int {
	result := 0

	if p.Kind == term {
		result++
	}

	result += strings.Count(p.Key, term) // XXX: this will break unless the Key is split properly

	return result
}

// PostingsList is a list of Posting objects. When used in an index,
// the list is sorted by Frequency and then DocID.
type PostingsList []Posting

func (p PostingsList) Len() int {
	return len(p)
}

func (p PostingsList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PostingsList) Less(i, j int) bool {
	if p[i].Frequency == p[j].Frequency {
		return p[i].DocID() < p[j].DocID()
	}
	return p[j].Frequency < p[i].Frequency // Large TF (term frequency) should come before small TF
}
