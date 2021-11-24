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

// Handler is an http.HandlerFunc that responds with just "Hello World!".
func Handler(index *InvertedIndex, store cache.Store) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		values, ok := request.URL.Query()["query"]

		if !ok || values[0] == "" {
			writeEmptyOutput(writer)
			return
		}

		key, found := index.Get(values[0])

		if !found {
			writeEmptyOutput(writer)
			return
		}

		item, exists, err := store.GetByKey(key)

		if err != nil {
			klog.Errorln(err)
			writeEmptyOutput(writer)
			return
		}

		if !exists {
			writeEmptyOutput(writer)
			return
		}

		pod := item.(*corev1.Pod)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(writer, fmt.Sprintf(`{"kind":"Pods","namespace":"%s","name":"%s"}`, pod.Namespace, pod.Name))
	}
}

func writeEmptyOutput(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(writer, "{}")
}
