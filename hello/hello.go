// Package hello is a TDD starting point. It'll be deleted later.
package hello

import (
	"io"
	"net/http"
)

func init() {
	http.HandleFunc("/v1/hello", Handler)
}

// Handler is an http.HandlerFunc that responds with just "Hello World!".
func Handler(writer http.ResponseWriter, request *http.Request) {
	io.WriteString(writer, "<html><body>Hello World!</body></html>")
}
