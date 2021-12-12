package index

// SearchFunc is a basic search function.
type SearchFunc func(query string) []Posting

// Searcher returns a search functor.
func Searcher(idx *Index) SearchFunc {
	return func(query string) []Posting {
		return idx.Get(query)
	}
}
