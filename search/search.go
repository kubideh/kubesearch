// Package search provides the API for searching for Kubernetes
// objects. Currently, just one method for querying exists, and
// it's endpoint is `/v1/search?query=`.
package search

import (
	"encoding/json"
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

// Result is a single result entry.
type Result struct {
	Kind      string `json:"kind,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespaces,omitempty"`
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

		objects, err := resultsFromPostings(postings, store)

		if err != nil {
			klog.Errorln(err)
			writeEmptyOutput(writer)
			return
		}

		writeResults(writer, objects)
	}
}

func writeResults(writer http.ResponseWriter, objects []Result) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	io.WriteString(writer, "[")

	for i, o := range objects {
		if i != 0 {
			io.WriteString(writer, ",")
		}
		entry, err := json.Marshal(o)
		if err != nil {
			klog.Warningln("error marshaling result: ", err)
			continue
		}
		io.WriteString(writer, string(entry))
	}

	io.WriteString(writer, "]")
}

func resultsFromPostings(postings []Posting, store map[string]cache.Store) ([]Result, error) {
	var results []Result

	for _, p := range postings {
		item, exists, err := store[p.Kind].GetByKey(p.Key)

		if err != nil {
			return results, err
		}

		if !exists {
			return results, err
		}

		// XXX Refactor this to make it dynamic; not dependent on Kind.

		switch p.Kind {
		case "Deployment":
			results = append(results, resultFromDeployment(item.(*appsv1.Deployment)))
		case "Pod":
			results = append(results, resultFromPod(item.(*corev1.Pod)))
		}
	}

	return results, nil
}

func resultFromDeployment(deployment *appsv1.Deployment) Result {
	return Result{
		Kind:      "Deployment",
		Name:      deployment.GetName(),
		Namespace: deployment.GetNamespace(),
	}
}

func resultFromPod(pod *corev1.Pod) Result {
	return Result{
		Kind:      "Pod",
		Name:      pod.GetName(),
		Namespace: pod.GetNamespace(),
	}
}

func writeEmptyOutput(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(writer, "[]")
}
