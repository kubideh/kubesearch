package searcher

import (
	"github.com/kubideh/kubesearch/search/index"
	"github.com/kubideh/kubesearch/search/tokenizer"
)

// SearchFunc is a basic search function.
type SearchFunc func(query string) []index.Posting

// Searcher returns the default search functor.
func Searcher(idx *index.Index, tokenize tokenizer.TokenizeFunc) SearchFunc {
	return func(query string) []index.Posting {
		terms := tokenize(query)
		//sort.Strings(terms)

		var result []index.Posting

		for _, t := range terms {
			postings := idx.Get(t)
			result = intersect(result, postings)
		}

		return result
	}
}

func intersect(left, right []index.Posting) (result []index.Posting) {

	if len(left) == 0 {
		return right
	}

	if len(right) == 0 {
		return left
	}

	i := 0
	j := 0

	for i < len(left) && j < len(right) {
		if left[i].DocID() == right[j].DocID() {
			result = append(result, postingWithLargestTermFrequency(left[i], right[j]))
			i++
			j++
		} else if left[i].DocID() < right[j].DocID() {
			i++
		} else {
			j++
		}
	}

	return
}

func postingWithLargestTermFrequency(p1, p2 index.Posting) index.Posting {
	if p1.Frequency < p2.Frequency {
		return p2
	}
	return p1
}
