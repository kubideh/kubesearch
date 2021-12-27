// Package api provides the API for searching for Kubernetes
// objects. Currently, just one method for querying exists.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kubideh/kubesearch/search/searcher"

	"github.com/kubideh/kubesearch/search/index"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

const (
	endpointPath   = "/v1/search"
	queryParamName = "query"
)

// RegisterHandler registers the search API handler with the given mux.
func RegisterHandler(mux *http.ServeMux, search searcher.SearchFunc, store map[string]cache.Store) {
	mux.HandleFunc(endpointPath, Handler(search, store))
}

// Handler is an http.HandlerFunc that responds with query results.
func Handler(search searcher.SearchFunc, store map[string]cache.Store) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		postings := search(query(request))

		objects, err := lookupObjects(postings, store)

		if err != nil {
			klog.Errorln(err)
		}

		writeResults(writer, objects)
	}
}

func query(request *http.Request) string {
	values, ok := request.URL.Query()[queryParamName]

	if !ok || len(values) == 0 {
		return ""
	}

	return values[0]
}

// Result is a single result entry.
type Result struct {
	Kind      string `json:"kind,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespaces,omitempty"`
	Rank      int    `json:"rank,omitempty"`
}

func writeResults(writer http.ResponseWriter, objects []Result) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	encoder := json.NewEncoder(writer)

	if err := encoder.Encode(objects); err != nil {
		klog.Warningln("error marshaling result: ", err)
	}
}

func lookupObjects(postings []index.Posting, store map[string]cache.Store) ([]Result, error) {
	var results []Result

	for _, p := range postings {
		item, exists, err := store[p.Kind].GetByKey(p.Key)

		if err != nil {
			return results, err
		}

		if !exists {
			return results, fmt.Errorf("missing object for key %v", p.Key)
		}

		// XXX Refactor this to make it dynamic; not dependent on Kind.

		switch p.Kind {
		case "Deployment":
			results = append(results, resultFromDeployment(item.(*appsv1.Deployment), p.Frequency))
		case "Pod":
			results = append(results, resultFromPod(item.(*corev1.Pod), p.Frequency))
		}
	}

	return results, nil
}

func resultFromDeployment(deployment *appsv1.Deployment, termFrequency int) Result {
	return Result{
		Kind:      "Deployment",
		Name:      deployment.GetName(),
		Namespace: deployment.GetNamespace(),
		Rank:      termFrequency,
	}
}

func resultFromPod(pod *corev1.Pod, termFrequency int) Result {
	return Result{
		Kind:      "Pod",
		Name:      pod.GetName(),
		Namespace: pod.GetNamespace(),
		Rank:      termFrequency,
	}
}
