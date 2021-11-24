package search

// InvertedIndex maps terms to object keys.
type InvertedIndex map[string]string

// NewIndex returns a Index objects.
func NewIndex() InvertedIndex {
	return make(InvertedIndex)
}
