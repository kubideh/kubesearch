//go:build integration
// +build integration

package hello

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloAPI(t *testing.T) {
	response, err := http.Get("http://localhost:8080/v1/search")

	require.NoError(t, err)

	body, _ := io.ReadAll(response.Body)
	response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, "<html><body>Hello World!</body></html>", string(body))
}
