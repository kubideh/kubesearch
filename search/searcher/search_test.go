package searcher

import (
	"testing"

	"github.com/kubideh/kubesearch/search/index"
	"github.com/kubideh/kubesearch/search/tokenizer"
	"github.com/stretchr/testify/assert"
)

func TestSearch_emptyQuery(t *testing.T) {
	search := Searcher(index.New(), tokenizer.Tokenizer())

	assert.Empty(t, search(""))
}

func TestSearch_missingObject(t *testing.T) {
	search := Searcher(index.New(), tokenizer.Tokenizer())

	assert.Empty(t, search("blargle"))
}

func TestSearch_singleTermMatchesOneObject(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{Key: "blargle", Kind: "flargle"})
	idx.Put([]string{"blargle"}, index.Posting{Key: "blargle", Kind: "flargle"})
	search := Searcher(idx, tokenizer.Tokenizer())

	assert.Equal(t, []index.Posting{{Key: "blargle", Kind: "flargle"}}, search("blargle"))
}

func TestSearch_singleTermMatchesTwoObjects(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{Key: "blargle", Kind: "flargle"})
	idx.Put([]string{"blargle"}, index.Posting{Key: "blargle", Kind: "bobble"})
	search := Searcher(idx, tokenizer.Tokenizer())

	assert.Equal(t, []index.Posting{
		{Key: "blargle", Kind: "bobble"},
		{Key: "blargle", Kind: "flargle"},
	}, search("blargle"))
}

func TestSearch_multipleTermsMatchTheSameObject(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})
	search := Searcher(idx, tokenizer.Tokenizer())

	assert.Equal(t, []index.Posting{{Key: "flargle/blargle", Kind: "flargle"}}, search("blargle flargle"))
}

func TestSearch_multipleTermsInDifferentOrderMatchTheSameObject(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})
	search := Searcher(idx, tokenizer.Tokenizer())

	assert.Equal(t, []index.Posting{{Key: "flargle/blargle", Kind: "flargle"}}, search("flargle blargle"))
}
