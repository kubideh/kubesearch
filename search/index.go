package search

var singletonIndex map[string]string

// Index returns the global pod search index.
func Index() map[string]string {
	return singletonIndex
}

// SetIndex replaces the global pod search index.
func SetIndex(index map[string]string) {
	singletonIndex = index
}
