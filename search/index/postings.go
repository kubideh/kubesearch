package index

import (
	"fmt"
	"strings"
)

// Posting represents a stored object key, what kind of K8s resource
// that object is, and a frequency of the number of times a
// particular term was found in that object.
type Posting struct {
	StoredObjectKey string
	K8sResourceKind string
	TermFrequency   int
}

// DocID wraps the document ID of a particular Posting, and it has
// the form <kind>/<namespace>/<objectID>.
type DocID struct {
	id string
}

func (d DocID) String() string {
	return d.id
}

// DocID is the document identifier, and it's a string with the
// form <K8sResourceKind>/<Optional namespace>/<Object name>.
func (p Posting) DocID() DocID {
	return DocID{id: fmt.Sprintf("%s/%s", p.K8sResourceKind, p.StoredObjectKey)}
}

// ComputeTermFrequency returns the number of times term appears in
// the given Posting.
func (p Posting) ComputeTermFrequency(term string) int {
	result := 0

	if p.K8sResourceKind == term {
		result++
	}

	result += strings.Count(p.StoredObjectKey, term) // XXX: this will break unless the StoredObjectKey is split properly

	return result
}

// PostingsList is a list of Posting objects. When used in an
// index, the list is sorted by largest TermFrequency and then
// DocID.
type PostingsList []Posting

func (p PostingsList) Len() int {
	return len(p)
}

func (p PostingsList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PostingsList) Less(i, j int) bool {
	if p[i].TermFrequency == p[j].TermFrequency {
		return p[i].DocID().String() < p[j].DocID().String()
	}
	return p[j].TermFrequency < p[i].TermFrequency // Large TF (term frequency) should come before small TF
}
