// Package search provides the API for searching for Kubernetes objects.
package search

import (
	"context"
	"fmt"
	"io"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var client kubernetes.Interface

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

	pod, err := client.CoreV1().Pods("flargle").Get(context.TODO(), keys[0], metav1.GetOptions{})

	if err != nil {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(writer, "{}")
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(writer, fmt.Sprintf(`{"kind":"Pods","namespace":"%s","name":%s"}`, pod.Namespace, pod.Name))
}
