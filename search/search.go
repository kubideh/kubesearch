// Package search provides the API for searching for Kubernetes
// objects. Currently, just one method for querying exists, and
// it's endpoint is `/v1/search?query=`.
package search

import (
	"fmt"
	"io"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

func init() {
	http.HandleFunc("/v1/search", Handler)
}

// Handler is an http.HandlerFunc that responds with just "Hello World!".
func Handler(writer http.ResponseWriter, request *http.Request) {
	keys, ok := request.URL.Query()["query"]

	if !ok || keys[0] == "" {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(writer, "{}")
		return
	}

	key, err := cache.MetaNamespaceKeyFunc(
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "flargle",
				Name:      keys[0],
			},
		},
	)

	if err != nil {
		klog.Error(err)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(writer, "{}")
		return
	}

	item, exists, err := Informer().GetStore().GetByKey(key)

	if err != nil {
		klog.Error(err)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(writer, "{}")
		return
	}

	if !exists {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(writer, "{}")
		return
	}

	pod := item.(*corev1.Pod)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(writer, fmt.Sprintf(`{"kind":"Pods","namespace":"%s","name":"%s"}`, pod.Namespace, pod.Name))
}
