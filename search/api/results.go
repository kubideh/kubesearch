package api

import (
	"github.com/kubideh/kubesearch/search/finder"
	"github.com/kubideh/kubesearch/search/index"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// Result is a single result entry.
type Result struct {
	Kind      string `json:"kind,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespaces,omitempty"`
	Rank      int    `json:"rank,omitempty"`
}

func createResults(objects []finder.Object, postings []index.Posting) (results []Result) {
	for i, o := range objects {
		results = append(results, createResult(postings[i].K8sResourceKind, o.Item, postings[i].TermFrequency))
	}
	return
}

func createResult(k8sResourceKind string, item interface{}, termFrequency int) (result Result) {
	switch k8sResourceKind {
	case "Deployment":
		result = createResultFromDeployment(item.(*appsv1.Deployment), termFrequency)
	case "Pod":
		result = createResultFromPod(item.(*corev1.Pod), termFrequency)
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
