package search

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

type testSearchCase struct {
	name   string
	params string
	result string
}

func TestSearch(t *testing.T) {
	cases := []testSearchCase{
		{
			name:   "search for pod by name",
			params: "query=blargle",
			result: `{"kind":"Pods","namespace":"flargle","name":"blargle"}`,
		},
		{
			name:   "search for missing object",
			params: "query=whatever",
			result: `{}`,
		},
		{
			name:   "search using empty query",
			params: "query=",
			result: `{}`,
		},
		{
			name:   "search with missing query param",
			params: "",
			result: `{}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { testSearch(t, c) })
	}
}

func testSearch(t *testing.T, c testSearchCase) {
	server, cancel := setup(t)
	defer server.Close()
	defer cancel()

	params := url.Values{}
	if c.params != "" {
		params.Add(strings.Split(c.params, "=")[0], strings.Split(c.params, "=")[1])
	}

	uri := fmt.Sprintf("%s/v1/search?%s", server.URL, params.Encode())
	t.Log(uri)
	response, err := http.Get(uri)
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, c.result, string(body))
}
