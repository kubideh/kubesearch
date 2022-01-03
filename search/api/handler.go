// Package api provides the API for searching for Kubernetes
// objects. Currently, just one method for querying exists.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/kubideh/kubesearch/search/finder"
	"github.com/kubideh/kubesearch/search/index"
	"github.com/kubideh/kubesearch/search/searcher"
	"k8s.io/klog/v2"
)

const (
	endpointPath   = "/v1/search"
	queryParamName = "queryString"
)

// RegisterSearchHandler registers the search API handler with the given mux
// at the appropriate endpoint path.
func RegisterSearchHandler(mux *http.ServeMux, handler http.HandlerFunc) {
	mux.HandleFunc(endpointPath, handler)
}

// CreateSearchHandler is a `http.HandlerFunc` that responds with a
// list of JSON-encoded results based on the given query string.
func CreateSearchHandler(search searcher.SearchFunc, findAll finder.FindAllFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		postings := search(queryString(request))

		keys := createKeysFromPostings(postings)
		objects, err := findAll(keys)

		if err != nil {
			klog.Errorln(err)
		}

		results := createResults(objects, postings)
		writeResults(writer, results)
	}
}

func queryString(request *http.Request) string {
	values, ok := request.URL.Query()[queryParamName]

	if !ok || len(values) == 0 {
		return ""
	}

	return values[0]
}

func createKeysFromPostings(postings []index.Posting) []finder.Key {
	keys := make([]finder.Key, 0, len(postings))

	for _, p := range postings {
		keys = appendPostingAsKey(keys, p)
	}

	return keys
}

func appendPostingAsKey(keys []finder.Key, posting index.Posting) (result []finder.Key) {
	key := createKeyFromPosting(posting)
	result = append(keys, key)
	return
}

func createKeyFromPosting(p index.Posting) finder.Key {
	return finder.Key{
		StoredObjectKey: p.StoredObjectKey,
		K8sResourceKind: p.K8sResourceKind,
	}
}

func writeResults(writer http.ResponseWriter, objects []Result) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	encoder := json.NewEncoder(writer)

	if err := encoder.Encode(objects); err != nil {
		klog.Warningln("error marshaling result: ", err)
	}
}
