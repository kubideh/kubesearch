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

// RegisterHandler registers the search API handler with the given mux.
func RegisterHandler(mux *http.ServeMux, search searcher.SearchFunc, findAll finder.FindAllFunc) {
	mux.HandleFunc(endpointPath, Handler(search, findAll))
}

// Handler is a `http.HandlerFunc` that responds with a list of
// JSON-encoded results based on the given query string.
func Handler(search searcher.SearchFunc, findAll finder.FindAllFunc) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		postings := search(queryString(request))

		keys := keysFromPostings(postings)
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

func keysFromPostings(postings []index.Posting) []finder.Key {
	keys := make([]finder.Key, 0, len(postings))

	for _, p := range postings {
		key := keyFromPosting(p)
		keys = append(keys, key)
	}

	return keys
}

func keyFromPosting(p index.Posting) finder.Key {
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
