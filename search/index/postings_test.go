package index

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTermFrequency(t *testing.T) {
	posting := Posting{
		Key:  "flargle/blargle",
		Kind: "flargle",
	}

	computed := posting.TermFrequency("flargle")

	assert.Equal(t, 2, computed)
}

func TestSortPostings(t *testing.T) {
	postings := []Posting{
		{
			Key:       "flargle/bobble",
			Kind:      "flargle",
			Frequency: 2,
		},
		{
			Key:       "flargle/blargle",
			Kind:      "flargle",
			Frequency: 2,
		},
		{
			Key:       "flargle/flargle",
			Kind:      "flargle",
			Frequency: 3,
		},
	}

	sort.Sort(PostingsList(postings))

	expected := []Posting{
		{
			Key:       "flargle/flargle",
			Kind:      "flargle",
			Frequency: 3,
		},
		{
			Key:       "flargle/blargle",
			Kind:      "flargle",
			Frequency: 2,
		},
		{
			Key:       "flargle/bobble",
			Kind:      "flargle",
			Frequency: 2,
		},
	}

	assert.Equal(t, expected, postings)
}
