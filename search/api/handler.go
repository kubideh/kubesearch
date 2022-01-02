// Package api provides the API for searching for Kubernetes
// objects. Currently, just one method for querying exists.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/kubideh/kubesearch/search/finder"
	"github.com/kubideh/kubesearch/search/searcher"

	"github.com/kubideh/kubesearch/search/index"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

// Handler is a `http.HandlerFunc` that responds with queryString
// results.
func Handler(search searcher.SearchFunc, findAll finder.FindAllFunc) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		postings := search(queryString(request))

		objects, err := findAll(postings)

		if err != nil {
			klog.Errorln(err)
		}

		results := transformObjectsIntoResults(objects)

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

// XXX everything below this line should be made into a data mapper.
func transformObjectsIntoResults(objects []finder.Object) (results []Result) {
	for _, o := range objects {
		results = append(results, createResult(o.Item, o.Posting))
	}
	return
}

func createResult(item interface{}, posting index.Posting) (result Result) {
	switch posting.K8sResourceKind {
	case "Deployment":
		result = createResultFromDeployment(item.(*appsv1.Deployment), posting.TermFrequency)
	case "Pod":
		result = createResultFromPod(item.(*corev1.Pod), posting.TermFrequency)
	}

	return
}

func createResultFromDeployment(deployment *appsv1.Deployment, termFrequency int) Result {
	return Result{
		Kind:      "Deployment",
		Name:      deployment.GetName(),
		Namespace: deployment.GetNamespace(),
		Rank:      termFrequency,
	}
}

func createResultFromPod(pod *corev1.Pod, termFrequency int) Result {
	return Result{
		Kind:      "Pod",
		Name:      pod.GetName(),
		Namespace: pod.GetNamespace(),
		Rank:      termFrequency,
	}
}
