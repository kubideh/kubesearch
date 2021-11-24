package search

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*httptest.Server, context.CancelFunc) {
	client := fake.NewSimpleClientset()
	controller := NewController(client)
	index := NewIndex()

	_, err := client.CoreV1().Pods("flargle").Create(
		context.TODO(),
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "blargle",
			},
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err)

	mux := http.NewServeMux()
	RegisterHandler(mux, index, controller.Store())
	server := httptest.NewServer(mux)

	return server, controller.Start(index)
}

func TestSearch_podByName(t *testing.T) {
	t.Parallel()

	server, cancel := setup(t)
	defer server.Close()
	defer cancel()

	params := url.Values{}
	params.Add("query", "blargle")

	response, err := http.Get(fmt.Sprintf("%s/v1/search?%s", server.URL, params.Encode()))
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, `{"kind":"Pods","namespace":"flargle","name":"blargle"}`, string(body))
}

func TestSearch_nonExistentPod(t *testing.T) {
	t.Parallel()

	server, cancel := setup(t)
	defer server.Close()
	defer cancel()

	params := url.Values{}
	params.Add("query", "whatever")

	response, err := http.Get(fmt.Sprintf("%s/v1/search?%s", server.URL, params.Encode()))
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, `{}`, string(body))
}

func TestSearch_missingQuery(t *testing.T) {
	t.Parallel()

	server, cancel := setup(t)
	defer server.Close()
	defer cancel()

	response, err := http.Get(fmt.Sprintf("%s/v1/search", server.URL))
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, `{}`, string(body))
}
