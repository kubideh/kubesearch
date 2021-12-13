package searcher

import (
	"testing"

	"github.com/kubideh/kubesearch/search/index"
	"github.com/stretchr/testify/assert"
)

func TestSearch_emptyQuery(t *testing.T) {
	idx := index.New()
	search := Searcher(idx)

	assert.Empty(t, search(""))
}

func TestSearch_missingObject(t *testing.T) {
	idx := index.New()
	search := Searcher(idx)

	assert.Empty(t, search("blargle"))
}

func TestSearch_singleTermMatchesOneObject(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{Key: "blargle", Kind: "flargle"})
	idx.Put([]string{"blargle"}, index.Posting{Key: "blargle", Kind: "flargle"})
	search := Searcher(idx)

	assert.Equal(t, []index.Posting{{Key: "blargle", Kind: "flargle"}}, search("blargle"))
}

func TestSearch_singleTermMatchesTwoObjects(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{Key: "blargle", Kind: "flargle"})
	idx.Put([]string{"blargle"}, index.Posting{Key: "blargle", Kind: "bobble"})
	search := Searcher(idx)

	assert.Equal(t, []index.Posting{
		{Key: "blargle", Kind: "bobble"},
		{Key: "blargle", Kind: "flargle"},
	}, search("blargle"))
}

func TestSearch_multipleTermsMatchTheSameObject(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})
	search := Searcher(idx)

	assert.Equal(t, []index.Posting{{Key: "flargle/blargle", Kind: "flargle"}}, search("blargle flargle"))
}

func TestSearch_multipleTermsInDifferentOrderMatchTheSameObject(t *testing.T) {
	idx := index.New()
	idx.Put([]string{"blargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})
	idx.Put([]string{"flargle"}, index.Posting{Key: "flargle/blargle", Kind: "flargle"})
	search := Searcher(idx)

	assert.Equal(t, []index.Posting{{Key: "flargle/blargle", Kind: "flargle"}}, search("flargle blargle"))
}
