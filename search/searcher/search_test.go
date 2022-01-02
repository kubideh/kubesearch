package searcher

import (
	"testing"

	"github.com/kubideh/kubesearch/search/index"
	"github.com/kubideh/kubesearch/search/tokenizer"
	"github.com/stretchr/testify/assert"
)

func TestSearch_emptyQuery(t *testing.T) {
	search := Create(index.New(), tokenizer.Tokenizer())

	result := search("")

	assert.Empty(t, result)
}

func TestSearch_missingObject(t *testing.T) {
	search := Create(index.New(), tokenizer.Tokenizer())

	result := search("blargle")

	assert.Empty(t, result)
}

func TestSearch_singleTermMatchesOneObject(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{StoredObjectKey: "blargle", K8sResourceKind: "flargle"})
	idx.Put([]string{"blargle"}, index.Posting{StoredObjectKey: "blargle", K8sResourceKind: "flargle"})

	search := Create(idx, tokenizer.Tokenizer())

	result := search("blargle")

	assert.Equal(t, []index.Posting{{StoredObjectKey: "blargle", K8sResourceKind: "flargle", TermFrequency: 1}}, result)
}

func TestSearch_singleTermMatchesTwoObjects(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{StoredObjectKey: "blargle", K8sResourceKind: "flargle"})
	idx.Put([]string{"blargle"}, index.Posting{StoredObjectKey: "blargle", K8sResourceKind: "bobble"})

	search := Create(idx, tokenizer.Tokenizer())

	result := search("blargle")

	assert.Equal(t, []index.Posting{
		{StoredObjectKey: "blargle", K8sResourceKind: "bobble", TermFrequency: 1},
		{StoredObjectKey: "blargle", K8sResourceKind: "flargle", TermFrequency: 1},
	}, result)
}

func TestSearch_multipleTermsMatchTheSameObject(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{StoredObjectKey: "flargle/blargle", K8sResourceKind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{StoredObjectKey: "flargle/blargle", K8sResourceKind: "flargle"})

	search := Create(idx, tokenizer.Tokenizer())

	result := search("blargle flargle")

	assert.Equal(t, []index.Posting{{StoredObjectKey: "flargle/blargle", K8sResourceKind: "flargle", TermFrequency: 2}}, result)
}

func TestSearch_multipleTermsInDifferentOrderMatchTheSameObject(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{StoredObjectKey: "flargle/blargle", K8sResourceKind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{StoredObjectKey: "flargle/blargle", K8sResourceKind: "flargle"})

	search := Create(idx, tokenizer.Tokenizer())

	result := search("flargle blargle")

	assert.Equal(t, []index.Posting{{StoredObjectKey: "flargle/blargle", K8sResourceKind: "flargle", TermFrequency: 2}}, result)
}

func TestSearch_orderedByRankAndDocID(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{StoredObjectKey: "flargle/blargle", K8sResourceKind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{StoredObjectKey: "flargle/blargle", K8sResourceKind: "flargle"})
	idx.Put([]string{"bobble"}, index.Posting{StoredObjectKey: "flargle/bobble", K8sResourceKind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{StoredObjectKey: "flargle/bobble", K8sResourceKind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{StoredObjectKey: "flargle/flargle", K8sResourceKind: "flargle"})

	search := Create(idx, tokenizer.Tokenizer())

	result := search("flargle")

	expected := []index.Posting{
		{StoredObjectKey: "flargle/flargle", K8sResourceKind: "flargle", TermFrequency: 3},
		{StoredObjectKey: "flargle/blargle", K8sResourceKind: "flargle", TermFrequency: 2},
		{StoredObjectKey: "flargle/bobble", K8sResourceKind: "flargle", TermFrequency: 2},
	}

	assert.Equal(t, expected, result)
}
