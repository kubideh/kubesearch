// Package search provides the API for searching for Kubernetes objects.
package search

import (
	"io"
	"net/http"
)

func init() {
	http.HandleFunc("/v1/search", Handler)
}

// Handler is an http.HandlerFunc that responds with just "Hello World!".
func Handler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(writer, "{}")
}
