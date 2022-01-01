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
	idx.Put([]string{"blargle"}, index.Posting{Key: "blargle", Kind: "flargle"})
	idx.Put([]string{"blargle"}, index.Posting{Key: "blargle", Kind: "flargle"})

	search := Create(idx, tokenizer.Tokenizer())

	result := search("blargle")

	assert.Equal(t, []index.Posting{{Key: "blargle", Kind: "flargle", Frequency: 1}}, result)
}

func TestSearch_singleTermMatchesTwoObjects(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{Key: "blargle", Kind: "flargle"})
	idx.Put([]string{"blargle"}, index.Posting{Key: "blargle", Kind: "bobble"})

	search := Create(idx, tokenizer.Tokenizer())

	result := search("blargle")

	assert.Equal(t, []index.Posting{
		{Key: "blargle", Kind: "bobble", Frequency: 1},
		{Key: "blargle", Kind: "flargle", Frequency: 1},
	}, result)
}

func TestSearch_multipleTermsMatchTheSameObject(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})

	search := Create(idx, tokenizer.Tokenizer())

	result := search("blargle flargle")

	assert.Equal(t, []index.Posting{{Key: "flargle/blargle", Kind: "flargle", Frequency: 2}}, result)
}

func TestSearch_multipleTermsInDifferentOrderMatchTheSameObject(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})

	search := Create(idx, tokenizer.Tokenizer())

	result := search("flargle blargle")

	assert.Equal(t, []index.Posting{{Key: "flargle/blargle", Kind: "flargle", Frequency: 2}}, result)
}

func TestSearch_orderedByRankAndDocID(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})
	idx.Put([]string{"bobble"}, index.Posting{Key: "flargle/bobble", Kind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{Key: "flargle/bobble", Kind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{Key: "flargle/flargle", Kind: "flargle"})

	search := Create(idx, tokenizer.Tokenizer())

	result := search("flargle")

	expected := []index.Posting{
		{Key: "flargle/flargle", Kind: "flargle", Frequency: 3},
		{Key: "flargle/blargle", Kind: "flargle", Frequency: 2},
		{Key: "flargle/bobble", Kind: "flargle", Frequency: 2},
	}

	assert.Equal(t, expected, result)
}
