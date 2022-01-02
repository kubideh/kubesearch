package index

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTermFrequency(t *testing.T) {
	posting := Posting{
		StoredObjectKey: "flargle/blargle",
		K8sResourceKind: "flargle",
	}

	computed := posting.ComputeTermFrequency("flargle")

	assert.Equal(t, 2, computed)
}

func TestSortPostings(t *testing.T) {
	postings := []Posting{
		{
			StoredObjectKey: "flargle/bobble",
			K8sResourceKind: "flargle",
			TermFrequency:   2,
		},
		{
			StoredObjectKey: "flargle/blargle",
			K8sResourceKind: "flargle",
			TermFrequency:   2,
		},
		{
			StoredObjectKey: "flargle/flargle",
			K8sResourceKind: "flargle",
			TermFrequency:   3,
		},
	}

	sort.Sort(PostingsList(postings))

	expected := []Posting{
		{
			StoredObjectKey: "flargle/flargle",
			K8sResourceKind: "flargle",
			TermFrequency:   3,
		},
		{
			StoredObjectKey: "flargle/blargle",
			K8sResourceKind: "flargle",
			TermFrequency:   2,
		},
		{
			StoredObjectKey: "flargle/bobble",
			K8sResourceKind: "flargle",
			TermFrequency:   2,
		},
	}

	assert.Equal(t, expected, postings)
}
