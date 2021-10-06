package hello

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHello(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://localhost/v1/search", nil)
	writer := httptest.NewRecorder()

	Handler(writer, request)

	response := writer.Result()
	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, "<html><body>Hello World!</body></html>", string(body))
}
