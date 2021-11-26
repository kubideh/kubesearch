// Package search provides the API for searching for Kubernetes
// objects. Currently, just one method for querying exists, and
// it's endpoint is `/v1/search?query=`.
package search

import (
	"fmt"
	"io"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

// RegisterHandler registers the search API handler with the given mux.
func RegisterHandler(mux *http.ServeMux, index *InvertedIndex, store cache.Store) {
	mux.HandleFunc("/v1/search", Handler(index, store))
}

// Handler is an http.HandlerFunc that responds with just "Hello World!".
func Handler(index *InvertedIndex, store cache.Store) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		values, ok := request.URL.Query()["query"]

		if !ok || values[0] == "" {
			writeEmptyOutput(writer)
			return
		}

		postings, found := index.Get(values[0])

		if !found {
			writeEmptyOutput(writer)
			return
		}

		var pods []*corev1.Pod
		for _, p := range postings {
			item, exists, err := store.GetByKey(p)

			if err != nil {
				klog.Errorln(err)
				writeEmptyOutput(writer)
				return
			}

			if !exists {
				writeEmptyOutput(writer)
				return
			}

			pods = append(pods, item.(*corev1.Pod))
		}

		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(writer, "[")
		for i, p := range pods {
			if i != 0 {
				io.WriteString(writer, ",")
			}
			io.WriteString(writer, fmt.Sprintf(`{"kind":"Pods","namespace":"%s","name":"%s"}`, p.Namespace, p.Name))
		}
		io.WriteString(writer, "]")
	}
}

func writeEmptyOutput(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(writer, "[]")
}
