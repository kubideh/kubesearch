// Package search provides the API for searching for Kubernetes
// objects. Currently, just one method for querying exists, and
// it's endpoint is `/v1/search?query=`.
package search

import (
	"fmt"
	"io"
	"net/http"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

// RegisterHandler registers the search API handler with the given mux.
func RegisterHandler(mux *http.ServeMux, index *Index, store map[string]cache.Store) {
	mux.HandleFunc("/v1/search", Handler(index, store))
}

// Handler is an http.HandlerFunc that responds with query results.
func Handler(index *Index, store map[string]cache.Store) func(http.ResponseWriter, *http.Request) {
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

		var objects []string
		for _, p := range postings {
			item, exists, err := store[p.Kind].GetByKey(p.Key)

			if err != nil {
				klog.Errorln(err)
				writeEmptyOutput(writer)
				return
			}

			if !exists {
				writeEmptyOutput(writer)
				return
			}

			// XXX Refactor this to make it dynamic; not dependent on Kind.

			switch p.Kind {
			case "Deployment":
				deployment := item.(*appsv1.Deployment)
				objects = append(objects, fmt.Sprintf(`{"kind":"Deployment","namespace":"%s","name":"%s"}`, deployment.Namespace, deployment.Name))
			case "Pod":
				pod := item.(*corev1.Pod)
				objects = append(objects, fmt.Sprintf(`{"kind":"Pod","namespace":"%s","name":"%s"}`, pod.Namespace, pod.Name))
			}
		}

		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(writer, "[")
		for i, o := range objects {
			if i != 0 {
				io.WriteString(writer, ",")
			}
			io.WriteString(writer, o)
		}
		io.WriteString(writer, "]")
	}
}

func writeEmptyOutput(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(writer, "[]")
}
